package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/go-faster/yaml"
	"github.com/ogen-go/ogen"
	"github.com/ogen-go/ogen/gen"
	"github.com/ogen-go/ogen/gen/ir"
	"go.uber.org/zap"
)

// analyze runs ogen's generator over the spec and derives
// the semantic Model from ogen's own intermediate representation. The
// returned Generator is the same instance used to write internal/ogen, so
// the analyzed IR and the generated client match by construction.
func analyze(specPath, ogenConfigPath string, cfg *Config) (*Model, *gen.Generator, error) {
	data, err := os.ReadFile(specPath)
	if err != nil {
		return nil, nil, err
	}
	spec, err := ogen.Parse(data)
	if err != nil {
		return nil, nil, fmt.Errorf("parse %s: %w", specPath, err)
	}
	opts, err := loadOgenOptions(ogenConfigPath)
	if err != nil {
		return nil, nil, err
	}
	g, err := gen.NewGenerator(spec, opts)
	if err != nil {
		return nil, nil, fmt.Errorf("ogen: %w", err)
	}

	l := &loader{
		cfg:      cfg,
		aliases:  newAliasRegistry(cfg.Rename.Aliases),
		optTypes: map[string]*ir.Type{},
	}
	model, err := l.build(g)
	if err != nil {
		return nil, nil, err
	}
	return model, g, nil
}

// loadOgenOptions reads the same config file the ogen CLI consumes;
// gen.Options carries the yaml tags for it.
func loadOgenOptions(path string) (gen.Options, error) {
	var opts gen.Options
	data, err := os.ReadFile(path)
	if err != nil {
		return opts, err
	}
	if err := yaml.Unmarshal(data, &opts); err != nil {
		return opts, fmt.Errorf("parse %s: %w", path, err)
	}
	opts.Logger = zap.NewNop()
	return opts, nil
}

type loader struct {
	cfg     *Config
	aliases *aliasRegistry
	// optTypes are ir types used in Options fields that need constructor
	// and enum re-exports (OptInt, OptListIncidentsSortOrder, ...).
	optTypes map[string]*ir.Type
}

func (l *loader) build(g *gen.Generator) (*Model, error) {
	// Fixed basic set, kept for compatibility with the former hand-written
	// types.go. Only aliased when the spec actually produces the type.
	types := g.Types()
	for _, base := range []string{"OptString", "OptDateTime", "OptInt", "OptNilInt"} {
		t, ok := types[base]
		if !ok {
			continue
		}
		if _, err := l.aliases.add(base, base); err != nil {
			return nil, err
		}
		l.optTypes[base] = t
	}

	ops := append([]*ir.Operation(nil), g.Operations()...)
	sort.Slice(ops, func(i, j int) bool { return ops[i].Name < ops[j].Name })

	services := map[string]*Service{}
	var problems []string
	for _, irop := range ops {
		id := irop.Spec.OperationID
		if id == "" {
			problems = append(problems, fmt.Sprintf("%s %s: missing operationId", irop.Spec.HTTPMethod, irop.Spec.Path))
			continue
		}
		if !l.cfg.wants(id) {
			continue
		}
		op, err := l.buildOp(irop)
		if err != nil {
			problems = append(problems, fmt.Sprintf("%s: %v", id, err))
			continue
		}
		tag, err := l.serviceTag(irop)
		if err != nil {
			problems = append(problems, fmt.Sprintf("%s: %v", id, err))
			continue
		}
		svc, ok := services[tag]
		if !ok {
			svc = &Service{Tag: tag}
			services[tag] = svc
		}
		svc.Ops = append(svc.Ops, op)
	}
	if len(problems) > 0 {
		return nil, fmt.Errorf("unsupported operations (fix the convention violation or add them to skip in .genwaroom.yaml):\n  %s",
			strings.Join(problems, "\n  "))
	}

	model := &Model{Types: l.typesFile()}
	for _, svc := range services {
		sort.Slice(svc.Ops, func(i, j int) bool { return svc.Ops[i].RawName < svc.Ops[j].RawName })
		model.Services = append(model.Services, svc)
	}
	sort.Slice(model.Services, func(i, j int) bool { return model.Services[i].Tag < model.Services[j].Tag })
	return model, nil
}

// serviceTag resolves the tag an operation is grouped under.
func (l *loader) serviceTag(irop *ir.Operation) (string, error) {
	if tag, ok := l.cfg.Groups[irop.Spec.OperationID]; ok {
		return tag, nil
	}
	if len(irop.Spec.Tags) > 0 {
		return irop.Spec.Tags[0], nil
	}
	// Fallback: first meaningful path segment.
	path := irop.Spec.Path.String()
	seg := strings.TrimPrefix(path, "/api/v0/")
	seg, _, _ = strings.Cut(seg, "/")
	if seg == "" {
		return "", fmt.Errorf("cannot determine group: no spec tag and unrecognized path %s", path)
	}
	tag := pascal(strings.ReplaceAll(seg, "_", " "))
	fmt.Fprintf(os.Stderr, "genwaroom: %s: no spec tag, grouping by path segment as %q\n", irop.Spec.OperationID, tag)
	return tag, nil
}

func (l *loader) buildOp(irop *ir.Operation) (*Op, error) {
	goName := irop.Name
	if forced, ok := l.cfg.Rename.Methods[irop.Spec.OperationID]; ok {
		goName = forced
	}
	op := &Op{
		ID:      irop.Spec.OperationID,
		GoName:  goName,
		RawName: irop.Name,
		Method:  strings.ToUpper(irop.Spec.HTTPMethod),
		Path:    irop.Spec.Path.String(),
		Summary: sentence(irop.Spec.Summary),
	}

	if irop.HasHeaderParams() {
		return nil, fmt.Errorf("header parameters are not supported")
	}
	if irop.HasCookieParams() {
		return nil, fmt.Errorf("cookie parameters are not supported")
	}

	// ogen synthesizes a <Name>Params struct for any op with parameters.
	if len(irop.Params) > 0 {
		op.ParamsGen = irop.Name + "Params"
	}
	// PathParts gives path parameters in true path order.
	for _, part := range irop.PathParts {
		if part.Param == nil {
			continue
		}
		op.PathParams = append(op.PathParams, PathParam{
			Type:     part.Param.Type.Go(),
			GenField: part.Param.Name,
		})
	}
	var optFields []Field
	for _, p := range irop.QueryParams() {
		typ, err := l.optionType(p.Type)
		if err != nil {
			return nil, fmt.Errorf("query param %s: %w", p.Name, err)
		}
		optFields = append(optFields, Field{Name: p.Name, Type: typ})
	}
	if len(optFields) > 0 {
		op.Options = &Options{Name: goName + "Options", Fields: optFields}
		if err := l.aliases.reserve(op.Options.Name); err != nil {
			return nil, err
		}
	}

	if irop.Request != nil {
		req, err := l.requestShape(op, irop.Request)
		if err != nil {
			return nil, err
		}
		op.Request = req
	}

	if err := l.responseShape(op, irop); err != nil {
		return nil, err
	}

	// Auto-pagination variant.
	op.All = op.Success.Kind == "list" && op.Options != nil && hasPageParams(op.Options)

	// GoName is stored only when it differs from RawName (rename config).
	if op.GoName == op.RawName {
		op.GoName = ""
	}
	return op, nil
}

// requestShape decides whether the request wrapper is unwrapped and
// registers the public alias of the parameter type.
func (l *loader) requestShape(op *Op, r *ir.Request) (*Request, error) {
	if r.EmptyBody != nil {
		return nil, fmt.Errorf("optional request body is not supported")
	}
	rt := r.Type
	if rt == nil || rt.Kind != ir.KindStruct {
		return nil, fmt.Errorf("request type %s is not a struct", r.GoType())
	}
	reqName := rt.Name
	if len(rt.Fields) == 1 {
		f := rt.Fields[0]
		if f.Type.Kind == ir.KindStruct && f.Type.Name == reqName+f.Name {
			public, err := l.aliases.add(op.GoName+"Request", f.Type.Name)
			if err != nil {
				return nil, err
			}
			return &Request{Public: public, Gen: reqName, UnwrapField: f.Name}, nil
		}
	}
	public, err := l.aliases.add(op.GoName+"Request", reqName)
	if err != nil {
		return nil, err
	}
	return &Request{Public: public, Gen: reqName}, nil
}

// responseShape classifies the responses of the operation: one success
// variant plus error variants keyed by status code.
func (l *loader) responseShape(op *Op, irop *ir.Operation) error {
	var successes []string
	seen := map[int]bool{}
	for _, info := range irop.ListResponseTypes(false) {
		if seen[info.StatusCode] {
			return fmt.Errorf("multiple response types for status %d (multiple content types?)", info.StatusCode)
		}
		seen[info.StatusCode] = true
		if info.Type == nil || info.Type.Name == "" {
			return fmt.Errorf("unnamed response type for status %d", info.StatusCode)
		}
		variant := info.Type.Name

		if info.StatusCode >= 200 && info.StatusCode < 300 {
			successes = append(successes, variant)
			if info.NoContent {
				op.Success = Success{Kind: "none", Variant: variant}
				continue
			}
			if err := l.successShape(op, info.Type); err != nil {
				return err
			}
			continue
		}

		// Error variant: validate the body shape and assert the name
		// suffix agrees with the spec status (render derives the status
		// from the suffix).
		if err := validateErrBody(info.Type); err != nil {
			return err
		}
		status, err := statusOf(variant, irop.Name)
		if err != nil {
			return err
		}
		if status != info.StatusCode {
			return fmt.Errorf("variant %s: name implies status %d but spec says %d", variant, status, info.StatusCode)
		}
		op.Errors = append(op.Errors, variant)
	}
	if len(successes) != 1 {
		return fmt.Errorf("expected exactly one success variant, got %v", successes)
	}
	sort.Strings(op.Errors)
	return nil
}

// validateErrBody asserts the error variant body is one of the three
// shapes newAPIError extracts messages from. The shape itself is not
// stored — the runtime detects it through ogen's getters — but an unknown
// shape must fail here instead of silently yielding empty messages.
func validateErrBody(t *ir.Type) error {
	if t.Kind != ir.KindStruct || len(t.Fields) != 1 {
		return fmt.Errorf("error variant %s: expected a struct with 1 field", t.Name)
	}
	f := t.Fields[0]
	switch {
	case f.Name == "Errors" && f.Type.Kind == ir.KindArray && f.Type.Item != nil && f.Type.Item.Go() == "string":
		return nil
	case f.Name == "Error" && f.Type.Name == "OptString":
		return nil
	case f.Name == "Error" && f.Type.Go() == "string":
		return nil
	}
	return fmt.Errorf("error variant %s has unknown body shape (%s %s)", t.Name, f.Name, f.Type.Go())
}

// successShape classifies the success variant into list / unwrap / whole
// and registers the public aliases it needs.
func (l *loader) successShape(op *Op, t *ir.Type) error {
	op.Success.Variant = t.Name
	if t.Kind != ir.KindStruct {
		return fmt.Errorf("success variant %s is not a struct", t.Name)
	}

	// list: exactly two fields, items slice + pagination metadata.
	if len(t.Fields) == 2 {
		var items, meta *ir.Field
		for _, f := range t.Fields {
			if f.Name == "ResponseMetadata" {
				meta = f
			} else if f.Type.Kind == ir.KindArray {
				items = f
			}
		}
		if items != nil && meta != nil {
			if err := checkPaginationMeta(meta.Type); err != nil {
				return fmt.Errorf("%s: %w", t.Name, err)
			}
			elem := items.Type.Item
			if elem == nil || elem.Name == "" {
				return fmt.Errorf("%s.%s is not a slice of a named type", t.Name, items.Name)
			}
			base := singular(items.Name)
			itemPublic, err := l.aliases.addPreferred(elem.Name, base, base+"ListItem")
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
			op.Success.ItemsField = items.Name
			return nil
		}
	}

	// unwrap: single field with a variant-derived struct type.
	if len(t.Fields) == 1 {
		f := t.Fields[0]
		if f.Type.Kind == ir.KindStruct && f.Type.Name == t.Name+f.Name {
			public, err := l.aliases.addPreferred(f.Type.Name, unwrapName(op.GoName, f.Name))
			if err != nil {
				return err
			}
			op.Success.Kind = "unwrap"
			op.Success.UnwrapField = f.Name
			op.Success.Public = public
			return nil
		}
	}

	// whole: return the variant struct itself under a public alias.
	public, err := l.aliases.addPreferred(t.Name, wholeName(op.GoName))
	if err != nil {
		return err
	}
	op.Success.Kind = "whole"
	op.Success.Public = public
	return nil
}

// checkPaginationMeta asserts the exact ResponseMetadata shape all list
// operations share.
func checkPaginationMeta(t *ir.Type) error {
	if t.Kind != ir.KindStruct || len(t.Fields) != len(paginationFields) {
		return fmt.Errorf("ResponseMetadata does not have the expected %d fields", len(paginationFields))
	}
	for i, want := range paginationFields {
		if t.Fields[i].Name != want {
			return fmt.Errorf("ResponseMetadata field %d is %s, want %s", i, t.Fields[i].Name, want)
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

// optionType registers the alias of an Options field type (and remembers
// it for constructor/enum re-export).
func (l *loader) optionType(t *ir.Type) (string, error) {
	if t.Name == "" {
		return "", fmt.Errorf("unnamed type %s", t.Go())
	}
	public, err := l.aliases.add(t.Name, t.Name)
	if err != nil {
		return "", err
	}
	l.optTypes[t.Name] = t
	return public, nil
}

// typesFile assembles aliases, constructor re-exports, and enum constant
// re-exports for types_gen.go.
func (l *loader) typesFile() *TypesFile {
	tf := &TypesFile{}

	optNames := make([]string, 0, len(l.optTypes))
	for n := range l.optTypes {
		optNames = append(optNames, n)
	}
	sort.Strings(optNames)

	for _, n := range optNames {
		t := l.optTypes[n]
		// ogen generates a New<Type> constructor for every generic
		// (Opt/Nil) type.
		if t.Kind != ir.KindGeneric || t.GenericOf == nil {
			continue
		}
		inner := t.GenericOf
		var paramStr string
		if inner.Kind == ir.KindEnum && inner.Name != "" {
			// Enum parameter: alias it and re-export its constants.
			public, err := l.aliases.add(inner.Name, inner.Name)
			if err != nil {
				continue
			}
			paramStr = public
			// EnumVariant.Name is already the full constant name.
			consts := make([]string, 0, len(inner.EnumVariants))
			for _, v := range inner.EnumVariants {
				consts = append(consts, v.Name)
			}
			sort.Strings(consts)
			tf.Enums = append(tf.Enums, Enum{Type: public, Consts: consts})
		} else {
			paramStr = inner.Go()
		}
		tf.News = append(tf.News, NewFunc{Param: paramStr, Type: n})
	}
	tf.Aliases = l.aliases.sortedAliases()
	return tf
}
