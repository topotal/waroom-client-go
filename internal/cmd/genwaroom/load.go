package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"os"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
)

var (
	invokesRe = regexp.MustCompile(`^\S+ invokes (\S+) operation\.$`)
	routeRe   = regexp.MustCompile(`^(GET|PUT|POST|DELETE|OPTIONS|HEAD|PATCH|TRACE) (/\S*)$`)
	placehRe  = regexp.MustCompile(`\{([^}]+)\}`)
)

type docInfo struct {
	opID    string
	method  string
	path    string
	summary string
}

type loader struct {
	pkg     *packages.Package
	scope   *types.Scope
	tags    map[string]string
	cfg     *Config
	aliases *aliasRegistry
	// optTypes are gen types used in Options fields that need constructor
	// and enum re-exports (OptInt, OptListIncidentsSortOrder, ...).
	optTypes map[string]bool
}

// load analyzes the compiled ogen package and derives the semantic Model.
func load(src string, tags map[string]string, cfg *Config) (*Model, error) {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo,
	}, src)
	if err != nil {
		return nil, err
	}
	if len(pkgs) != 1 {
		return nil, fmt.Errorf("expected exactly one package for %s, got %d", src, len(pkgs))
	}
	pkg := pkgs[0]
	if len(pkg.Errors) > 0 {
		return nil, fmt.Errorf("load %s: %v", src, pkg.Errors[0])
	}

	l := &loader{
		pkg:      pkg,
		scope:    pkg.Types.Scope(),
		tags:     tags,
		cfg:      cfg,
		aliases:  newAliasRegistry(cfg.Rename.Aliases),
		optTypes: map[string]bool{},
	}
	return l.build()
}

func (l *loader) build() (*Model, error) {
	docs, err := l.invokerDocs()
	if err != nil {
		return nil, err
	}

	// Fixed basic set, kept for compatibility with the former hand-written
	// types.go.
	for _, base := range []string{"OptString", "OptDateTime", "OptInt", "OptNilInt"} {
		if _, err := l.aliases.add(base, base); err != nil {
			return nil, err
		}
		l.optTypes[base] = true
	}

	names := make([]string, 0, len(docs))
	for name := range docs {
		names = append(names, name)
	}
	sort.Strings(names)

	services := map[string]*Service{}
	var problems []string
	for _, name := range names {
		doc := docs[name]
		if !l.cfg.wants(doc.opID) {
			continue
		}
		op, err := l.buildOp(name, doc)
		if err != nil {
			problems = append(problems, fmt.Sprintf("%s: %v", doc.opID, err))
			continue
		}
		svc, err := l.serviceFor(doc)
		if err != nil {
			problems = append(problems, fmt.Sprintf("%s: %v", doc.opID, err))
			continue
		}
		if existing, ok := services[svc.Name]; ok {
			svc = existing
		} else {
			services[svc.Name] = svc
		}
		svc.Ops = append(svc.Ops, op)
	}
	if len(problems) > 0 {
		return nil, fmt.Errorf("unsupported operations (fix the convention violation or add them to skip in .genwaroom.yaml):\n  %s",
			strings.Join(problems, "\n  "))
	}

	model := &Model{Types: l.typesFile()}
	for _, svc := range services {
		sort.Slice(svc.Ops, func(i, j int) bool { return svc.Ops[i].GoName < svc.Ops[j].GoName })
		model.Services = append(model.Services, svc)
	}
	sort.Slice(model.Services, func(i, j int) bool { return model.Services[i].Name < model.Services[j].Name })
	return model, nil
}

// invokerDocs parses the Invoker interface's per-method doc comments.
func (l *loader) invokerDocs() (map[string]docInfo, error) {
	for _, f := range l.pkg.Syntax {
		for _, decl := range f.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.TYPE {
				continue
			}
			for _, spec := range gd.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok || ts.Name.Name != "Invoker" {
					continue
				}
				it, ok := ts.Type.(*ast.InterfaceType)
				if !ok {
					return nil, fmt.Errorf("Invoker is not an interface")
				}
				docs := map[string]docInfo{}
				for _, m := range it.Methods.List {
					if len(m.Names) == 0 {
						continue
					}
					name := m.Names[0].Name
					info, err := parseDoc(m.Doc)
					if err != nil {
						return nil, fmt.Errorf("Invoker.%s: %w", name, err)
					}
					docs[name] = info
				}
				return docs, nil
			}
		}
	}
	return nil, fmt.Errorf("Invoker interface not found in %s", l.pkg.PkgPath)
}

func parseDoc(cg *ast.CommentGroup) (docInfo, error) {
	var info docInfo
	if cg == nil {
		return info, fmt.Errorf("missing doc comment")
	}
	var body []string
	for _, line := range strings.Split(strings.TrimSpace(cg.Text()), "\n") {
		line = strings.TrimSpace(line)
		if m := invokesRe.FindStringSubmatch(line); m != nil {
			info.opID = m[1]
			continue
		}
		if m := routeRe.FindStringSubmatch(line); m != nil {
			info.method, info.path = m[1], m[2]
			continue
		}
		if line != "" {
			body = append(body, line)
		}
	}
	if info.opID == "" {
		return info, fmt.Errorf("doc comment lacks %q line", "<Name> invokes <operationId> operation.")
	}
	if info.method == "" {
		return info, fmt.Errorf("doc comment lacks %q line", "<METHOD> <path>")
	}
	info.summary = strings.Join(body, " ")
	return info, nil
}

// serviceFor resolves the service an operation belongs to.
func (l *loader) serviceFor(doc docInfo) (*Service, error) {
	tag, ok := l.cfg.Groups[doc.opID]
	if !ok {
		tag, ok = l.tags[doc.opID]
	}
	if !ok {
		// Fallback: first meaningful path segment.
		seg := strings.TrimPrefix(doc.path, "/api/v0/")
		seg, _, _ = strings.Cut(seg, "/")
		if seg == "" {
			return nil, fmt.Errorf("cannot determine group: no spec tag and unrecognized path %s", doc.path)
		}
		tag = pascal(strings.ReplaceAll(seg, "_", " "))
		fmt.Fprintf(os.Stderr, "genwaroom: %s: no spec tag, grouping by path segment as %q\n", doc.opID, tag)
	}
	base := pascal(tag)
	return &Service{
		Tag:      tag,
		Name:     singular(base) + "Service",
		Accessor: base,
		FileName: snake(base) + "_gen.go",
	}, nil
}

func (l *loader) buildOp(name string, doc docInfo) (*Op, error) {
	obj := l.scope.Lookup("Invoker")
	iface := obj.Type().Underlying().(*types.Interface)
	var sig *types.Signature
	for i := 0; i < iface.NumMethods(); i++ {
		if iface.Method(i).Name() == name {
			sig = iface.Method(i).Type().(*types.Signature)
			break
		}
	}
	if sig == nil {
		return nil, fmt.Errorf("method %s not found on Invoker", name)
	}

	goName := name
	if forced, ok := l.cfg.Rename.Methods[doc.opID]; ok {
		goName = forced
	}
	op := &Op{
		ID:      doc.opID,
		GoName:  goName,
		RawName: name,
		Method:  doc.method,
		Path:    doc.path,
		Summary: doc.summary,
	}

	reqNamed, paramsNamed, err := splitSignature(sig)
	if err != nil {
		return nil, err
	}

	// Path args and options from the params struct.
	if paramsNamed != nil {
		op.ParamsGen = paramsNamed.Obj().Name()
		st := paramsNamed.Underlying().(*types.Struct)
		pathParams, used, err := l.pathParams(doc.path, st)
		if err != nil {
			return nil, err
		}
		op.PathParams = pathParams

		var optFields []Field
		for i := 0; i < st.NumFields(); i++ {
			f := st.Field(i)
			if used[f.Name()] {
				continue
			}
			typ, err := l.optionType(f.Type())
			if err != nil {
				return nil, fmt.Errorf("query param %s: %w", f.Name(), err)
			}
			optFields = append(optFields, Field{Name: f.Name(), Type: typ})
		}
		if len(optFields) > 0 {
			op.Options = &Options{Name: goName + "Options", Fields: optFields}
			if err := l.aliases.reserve(op.Options.Name); err != nil {
				return nil, err
			}
		}
	}

	// Request body.
	if reqNamed != nil {
		req, err := l.requestShape(op, reqNamed)
		if err != nil {
			return nil, err
		}
		op.Request = req
	}

	// Response sum type.
	resNamed, ok := sig.Results().At(0).Type().(*types.Named)
	if !ok {
		return nil, fmt.Errorf("first result of %s is not a named type", name)
	}
	if err := l.responseShape(op, name, resNamed); err != nil {
		return nil, err
	}

	// Auto-pagination variant.
	op.All = op.Success.Kind == "list" && op.Options != nil && hasPageParams(op.Options)
	return op, nil
}

// splitSignature validates the uniform Invoker signature shape:
// (ctx, [request *XxxReq], [params XxxParams]) (XxxRes, error).
func splitSignature(sig *types.Signature) (req, params *types.Named, _ error) {
	tuple := sig.Params()
	if tuple.Len() == 0 || tuple.At(0).Type().String() != "context.Context" {
		return nil, nil, fmt.Errorf("first parameter is not context.Context")
	}
	for i := 1; i < tuple.Len(); i++ {
		t := tuple.At(i).Type()
		if ptr, ok := t.(*types.Pointer); ok {
			named, ok := ptr.Elem().(*types.Named)
			if !ok || req != nil || params != nil {
				return nil, nil, fmt.Errorf("unexpected parameter %s", t)
			}
			req = named
			continue
		}
		named, ok := t.(*types.Named)
		if !ok || params != nil {
			return nil, nil, fmt.Errorf("unexpected parameter %s", t)
		}
		params = named
	}
	if sig.Results().Len() != 2 {
		return nil, nil, fmt.Errorf("expected (Res, error) results")
	}
	return req, params, nil
}

// pathParams maps every {placeholder} in the path to a params-struct
// field, in path order. It also returns the set of consumed field names.
func (l *loader) pathParams(path string, st *types.Struct) ([]PathParam, map[string]bool, error) {
	var params []PathParam
	used := map[string]bool{}
	for _, m := range placehRe.FindAllStringSubmatch(path, -1) {
		want := normalizeParam(m[1])
		var field *types.Var
		for i := 0; i < st.NumFields(); i++ {
			if normalizeParam(st.Field(i).Name()) == want {
				field = st.Field(i)
				break
			}
		}
		if field == nil {
			return nil, nil, fmt.Errorf("path placeholder {%s} has no matching params field", m[1])
		}
		typ, err := l.externalType(field.Type())
		if err != nil {
			return nil, nil, fmt.Errorf("path param %s: %w", field.Name(), err)
		}
		params = append(params, PathParam{
			Arg:      lowerCamel(field.Name()),
			Type:     typ,
			GenField: field.Name(),
		})
		used[field.Name()] = true
	}
	return params, used, nil
}

func normalizeParam(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, "_", ""))
}

// requestShape decides whether the request wrapper is unwrapped and
// registers the public alias of the parameter type.
func (l *loader) requestShape(op *Op, reqNamed *types.Named) (*Request, error) {
	reqName := reqNamed.Obj().Name()
	st, ok := reqNamed.Underlying().(*types.Struct)
	if !ok {
		return nil, fmt.Errorf("request %s is not a struct", reqName)
	}
	if st.NumFields() == 1 {
		f := st.Field(0)
		if inner, ok := f.Type().(*types.Named); ok {
			if _, isStruct := inner.Underlying().(*types.Struct); isStruct && inner.Obj().Name() == reqName+f.Name() {
				public, err := l.aliases.add(op.GoName+"Request", inner.Obj().Name())
				if err != nil {
					return nil, err
				}
				return &Request{Public: public, Gen: reqName, UnwrapField: f.Name()}, nil
			}
		}
	}
	public, err := l.aliases.add(op.GoName+"Request", reqName)
	if err != nil {
		return nil, err
	}
	return &Request{Public: public, Gen: reqName}, nil
}

// responseShape classifies the success variant and error variants of the
// operation's Res sum type.
func (l *loader) responseShape(op *Op, rawName string, resNamed *types.Named) error {
	iface, ok := resNamed.Underlying().(*types.Interface)
	if !ok {
		return fmt.Errorf("%s is not an interface", resNamed.Obj().Name())
	}

	var successes []string
	for _, variant := range l.variants(iface) {
		suffix := strings.TrimPrefix(variant, rawName)
		if suffix == variant {
			return fmt.Errorf("response variant %s does not start with %s", variant, rawName)
		}
		if status, ok := errorStatus[suffix]; ok {
			body, err := l.errBody(variant)
			if err != nil {
				return err
			}
			op.Errors = append(op.Errors, ErrCase{Variant: variant, Status: status, Body: body})
			continue
		}
		if isSuccessSuffix(suffix) {
			successes = append(successes, variant)
			continue
		}
		return fmt.Errorf("response variant %s has unknown suffix %q", variant, suffix)
	}
	if len(successes) != 1 {
		return fmt.Errorf("expected exactly one success variant, got %v", successes)
	}
	return l.successShape(op, successes[0])
}

func isSuccessSuffix(suffix string) bool {
	for _, s := range successSuffixes {
		if s == suffix {
			return true
		}
	}
	return false
}

// variants returns all named struct types implementing the sealed Res
// interface, sorted by name.
func (l *loader) variants(iface *types.Interface) []string {
	var names []string
	for _, n := range l.scope.Names() {
		tn, ok := l.scope.Lookup(n).(*types.TypeName)
		if !ok || tn.IsAlias() {
			continue
		}
		named, ok := tn.Type().(*types.Named)
		if !ok {
			continue
		}
		if _, isIface := named.Underlying().(*types.Interface); isIface {
			continue
		}
		if types.Implements(types.NewPointer(named), iface) {
			names = append(names, n)
		}
	}
	return names
}

// errBody classifies the error variant body into one of the three known
// shapes.
func (l *loader) errBody(variant string) (string, error) {
	st, err := l.structOf(variant)
	if err != nil {
		return "", err
	}
	if st.NumFields() != 1 {
		return "", fmt.Errorf("error variant %s: expected 1 field, got %d", variant, st.NumFields())
	}
	f := st.Field(0)
	switch {
	case f.Name() == "Errors" && f.Type().String() == "[]string":
		return bodyErrorsArray, nil
	case f.Name() == "Error" && isGenNamed(f.Type(), "OptString"):
		return bodyErrorOpt, nil
	case f.Name() == "Error" && f.Type().String() == "string":
		return bodyErrorString, nil
	}
	return "", fmt.Errorf("error variant %s has unknown body shape (%s %s)", variant, f.Name(), f.Type())
}

// successShape classifies the success variant into list / unwrap / whole /
// none and registers the public aliases it needs.
func (l *loader) successShape(op *Op, variant string) error {
	op.Success.Variant = variant
	if strings.HasSuffix(variant, "NoContent") {
		op.Success.Kind = "none"
		return nil
	}
	st, err := l.structOf(variant)
	if err != nil {
		return err
	}

	// list: exactly two fields, items slice + pagination metadata.
	if st.NumFields() == 2 {
		itemsIdx, metaIdx := -1, -1
		for i := 0; i < st.NumFields(); i++ {
			if st.Field(i).Name() == "ResponseMetadata" {
				metaIdx = i
			} else if _, ok := st.Field(i).Type().(*types.Slice); ok {
				itemsIdx = i
			}
		}
		if itemsIdx >= 0 && metaIdx >= 0 {
			if err := checkPaginationMeta(st.Field(metaIdx).Type()); err != nil {
				return fmt.Errorf("%s: %w", variant, err)
			}
			items := st.Field(itemsIdx)
			elem, ok := items.Type().(*types.Slice).Elem().(*types.Named)
			if !ok {
				return fmt.Errorf("%s.%s is not a slice of a named type", variant, items.Name())
			}
			base := singular(items.Name())
			itemPublic, err := l.aliases.addPreferred(elem.Obj().Name(), base, base+"ListItem")
			if err != nil {
				return err
			}
			listName := base + "List"
			if err := l.aliases.reserve(listName); err != nil {
				return err
			}
			op.Success.Kind = "list"
			op.Success.ListName = listName
			op.Success.ItemPublic = itemPublic
			op.Success.ItemsField = items.Name()
			return nil
		}
	}

	// unwrap: single field with a variant-derived struct type.
	if st.NumFields() == 1 {
		f := st.Field(0)
		if inner, ok := f.Type().(*types.Named); ok {
			if _, isStruct := inner.Underlying().(*types.Struct); isStruct && inner.Obj().Name() == variant+f.Name() {
				public, err := l.aliases.addPreferred(inner.Obj().Name(), unwrapName(op.GoName, f.Name()))
				if err != nil {
					return err
				}
				op.Success.Kind = "unwrap"
				op.Success.UnwrapField = f.Name()
				op.Success.Public = public
				return nil
			}
		}
	}

	// whole: return the OK struct itself under a public alias.
	public, err := l.aliases.addPreferred(variant, wholeName(op.GoName))
	if err != nil {
		return err
	}
	op.Success.Kind = "whole"
	op.Success.Public = public
	return nil
}

// unwrapName derives the public name of an unwrapped response type:
// CreateServiceLabel + Label → CreatedLabel.
func unwrapName(goName, field string) string {
	verb := leadingWord(goName)
	if participle, ok := unwrapVerbs[verb]; ok {
		return participle + field
	}
	return goName + field
}

// wholeName derives the public name of a whole response type:
// GetIncident → Incident.
func wholeName(goName string) string {
	if rest := strings.TrimPrefix(goName, "Get"); rest != goName && rest != "" {
		return rest
	}
	return goName + "Result"
}

func leadingWord(s string) string {
	for i := 1; i < len(s); i++ {
		if s[i] >= 'A' && s[i] <= 'Z' {
			return s[:i]
		}
	}
	return s
}

// checkPaginationMeta asserts the exact ResponseMetadata shape all list
// operations share.
func checkPaginationMeta(t types.Type) error {
	named, ok := t.(*types.Named)
	if !ok {
		return fmt.Errorf("ResponseMetadata is not a named type")
	}
	st, ok := named.Underlying().(*types.Struct)
	if !ok || st.NumFields() != len(paginationFields) {
		return fmt.Errorf("ResponseMetadata does not have the expected %d fields", len(paginationFields))
	}
	for i, want := range paginationFields {
		if st.Field(i).Name() != want {
			return fmt.Errorf("ResponseMetadata field %d is %s, want %s", i, st.Field(i).Name(), want)
		}
	}
	return nil
}

func hasPageParams(o *Options) bool {
	var page, perPage bool
	for _, f := range o.Fields {
		if f.Name == "Page" && f.Type == "OptInt" {
			page = true
		}
		if f.Name == "PerPage" && f.Type == "OptInt" {
			perPage = true
		}
	}
	return page && perPage
}

// externalType renders a path-arg type (basic or external like uuid.UUID).
func (l *loader) externalType(t types.Type) (string, error) {
	switch t := t.(type) {
	case *types.Basic:
		return t.String(), nil
	case *types.Named:
		pkg := t.Obj().Pkg()
		if pkg == nil {
			return t.Obj().Name(), nil
		}
		if pkg == l.pkg.Types {
			return "", fmt.Errorf("gen-internal type %s not supported as a path argument", t.Obj().Name())
		}
		return pkg.Name() + "." + t.Obj().Name(), nil
	}
	return "", fmt.Errorf("unsupported type %s", t)
}

// optionType renders an Options field type, registering the alias (and
// remembering it for constructor/enum re-export).
func (l *loader) optionType(t types.Type) (string, error) {
	named, ok := t.(*types.Named)
	if !ok || named.Obj().Pkg() != l.pkg.Types {
		return "", fmt.Errorf("expected a gen-package named type, got %s", t)
	}
	name := named.Obj().Name()
	public, err := l.aliases.add(name, name)
	if err != nil {
		return "", err
	}
	l.optTypes[name] = true
	return public, nil
}

// typesFile assembles aliases, constructor re-exports, and enum constant
// re-exports for types_gen.go.
func (l *loader) typesFile() *TypesFile {
	tf := &TypesFile{Aliases: l.aliases.sortedAliases()}

	optNames := make([]string, 0, len(l.optTypes))
	for n := range l.optTypes {
		optNames = append(optNames, n)
	}
	sort.Strings(optNames)

	for _, n := range optNames {
		fn, ok := l.scope.Lookup("New" + n).(*types.Func)
		if !ok {
			continue
		}
		sig := fn.Type().(*types.Signature)
		if sig.Params().Len() != 1 {
			continue
		}
		param := sig.Params().At(0).Type()
		paramStr := ""
		switch p := param.(type) {
		case *types.Basic:
			paramStr = p.String()
		case *types.Named:
			if p.Obj().Pkg() == l.pkg.Types {
				// Enum (or other named) parameter: alias it and re-export
				// its constants.
				public, err := l.aliases.add(p.Obj().Name(), p.Obj().Name())
				if err != nil {
					continue
				}
				paramStr = public
				if consts := l.constsOf(p); len(consts) > 0 {
					tf.Enums = append(tf.Enums, Enum{Type: public, Consts: consts})
				}
			} else if p.Obj().Pkg() != nil {
				paramStr = p.Obj().Pkg().Name() + "." + p.Obj().Name()
			}
		}
		if paramStr == "" {
			continue
		}
		tf.News = append(tf.News, NewFunc{Name: "New" + n, Param: paramStr, Type: n})
	}
	// Aliases may have grown while resolving constructor params.
	tf.Aliases = l.aliases.sortedAliases()
	return tf
}

// constsOf lists package-level constants of the given named type.
func (l *loader) constsOf(named *types.Named) []string {
	var consts []string
	for _, n := range l.scope.Names() {
		c, ok := l.scope.Lookup(n).(*types.Const)
		if ok && c.Type() == named {
			consts = append(consts, n)
		}
	}
	sort.Strings(consts)
	return consts
}

func (l *loader) structOf(name string) (*types.Struct, error) {
	obj := l.scope.Lookup(name)
	if obj == nil {
		return nil, fmt.Errorf("type %s not found", name)
	}
	st, ok := obj.Type().Underlying().(*types.Struct)
	if !ok {
		return nil, fmt.Errorf("type %s is not a struct", name)
	}
	return st, nil
}

func isGenNamed(t types.Type, name string) bool {
	named, ok := t.(*types.Named)
	return ok && named.Obj().Name() == name
}
