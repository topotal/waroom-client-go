package main

import (
	"fmt"
	"sort"
	"strings"
)

// Model is the semantic intermediate representation. It is the contract
// between the two generator stages: `analyze` derives it from the ogen
// package and serializes it to YAML; `render` reads the YAML and emits Go
// code. It contains no Go source fragments — only names and shapes.
type Model struct {
	Services []*Service `yaml:"services"`
	Types    *TypesFile `yaml:"types"`
}

// Service groups operations by their spec tag.
type Service struct {
	Tag      string `yaml:"tag"`      // "Labels"
	Name     string `yaml:"name"`     // "LabelService"
	Accessor string `yaml:"accessor"` // "Labels"
	FileName string `yaml:"file"`     // "labels_gen.go"
	Ops      []*Op  `yaml:"operations"`
}

// Op is one generated public method (plus its optional Options struct and
// All iterator).
type Op struct {
	ID      string `yaml:"id"`      // operationId: createServiceLabel
	GoName  string `yaml:"goName"`  // public method name (usually = RawName)
	RawName string `yaml:"rawName"` // gen client method name
	Method  string `yaml:"method"`  // POST
	Path    string `yaml:"path"`    // /api/v0/services/{name}/labels
	Summary string `yaml:"summary,omitempty"`

	PathParams []PathParam `yaml:"pathParams,omitempty"` // in path order
	ParamsGen  string      `yaml:"paramsGen,omitempty"`  // gen params struct: CreateServiceLabelParams
	Request    *Request    `yaml:"request,omitempty"`
	Options    *Options    `yaml:"options,omitempty"`
	Success    Success     `yaml:"success"`
	Errors     []ErrCase   `yaml:"errors"`
	All        bool        `yaml:"all,omitempty"` // generate the auto-pagination iterator
}

// PathParam is a path placeholder promoted to a positional method argument.
type PathParam struct {
	Arg      string `yaml:"arg"`      // name
	Type     string `yaml:"type"`     // string / uuid.UUID
	GenField string `yaml:"genField"` // field in the gen params struct: Name
}

// Request describes the request body parameter.
type Request struct {
	Public      string `yaml:"public"`                // CreateServiceLabelRequest
	Gen         string `yaml:"gen"`                   // CreateServiceLabelReq
	UnwrapField string `yaml:"unwrapField,omitempty"` // Label; empty = pass the Req type as-is
}

// Success describes how the success variant maps to the public return.
type Success struct {
	Kind    string `yaml:"kind"`    // "list" | "unwrap" | "whole" | "none"
	Variant string `yaml:"variant"` // gen type name: ListIncidentsOK
	// kind=list
	ListName   string `yaml:"listName,omitempty"`   // IncidentList
	ItemPublic string `yaml:"itemPublic,omitempty"` // IncidentListItem
	ItemsField string `yaml:"itemsField,omitempty"` // Incidents
	// kind=unwrap
	UnwrapField string `yaml:"unwrapField,omitempty"` // Label
	// kind=unwrap|whole: public name of the returned type
	Public string `yaml:"public,omitempty"` // CreatedLabel / Incident
}

// ErrCase maps one error variant to an APIError.
type ErrCase struct {
	Variant string `yaml:"variant"` // gen type name: CreateServiceLabelUnauthorized
	Status  int    `yaml:"status"`
	Body    string `yaml:"body"` // bodyErrorsArray | bodyErrorOpt | bodyErrorString
}

// Error body shapes the conventions recognize.
const (
	bodyErrorsArray = "errorsArray" // Errors []string
	bodyErrorOpt    = "errorOpt"    // Error OptString
	bodyErrorString = "errorString" // Error string
)

// Options is the generated struct holding an op's query parameters.
type Options struct {
	Name   string  `yaml:"name"` // ListIncidentsOptions
	Fields []Field `yaml:"fields"`
}

// Field is one Options field; its name matches the gen params struct
// field, its type is a public alias name.
type Field struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"` // e.g. OptInt
}

// TypesFile drives types_gen.go.
type TypesFile struct {
	Aliases []Alias   `yaml:"aliases"`
	News    []NewFunc `yaml:"constructors,omitempty"`
	Enums   []Enum    `yaml:"enums,omitempty"`
}

// Alias is one `type Public = gen.Gen` line.
type Alias struct {
	Public string `yaml:"public"`
	Gen    string `yaml:"gen"`
}

// NewFunc re-exports an Opt/Nil constructor, e.g. NewOptInt.
type NewFunc struct {
	Name  string `yaml:"name"`  // NewOptInt
	Param string `yaml:"param"` // public param type: int or ListIncidentsSortOrder
	Type  string `yaml:"type"`  // public result type: OptInt
}

// Enum re-exports the constants of an aliased enum type.
type Enum struct {
	Type   string   `yaml:"type"` // ListIncidentsSortOrder
	Consts []string `yaml:"consts"`
}

// errorStatus maps ogen variant name suffixes to HTTP status codes. This
// is a closed set: an unknown suffix aborts generation.
var errorStatus = map[string]int{
	"BadRequest":          400,
	"Unauthorized":        401,
	"PaymentRequired":     402,
	"NotFound":            404,
	"UnprocessableEntity": 422,
}

var successSuffixes = []string{"OK", "Created", "NoContent"}

// paginationFields is the exact shape a list ResponseMetadata must have.
var paginationFields = []string{"CurrentPage", "NextPage", "PrevPage", "TotalPages", "TotalCount"}

// unwrapVerbs maps operation verbs to the participle used when naming an
// unwrapped response type: createServiceLabel + Label → CreatedLabel.
var unwrapVerbs = map[string]string{
	"Create": "Created",
	"Update": "Updated",
}

// pascal converts a tag like "Action Items" to "ActionItems".
func pascal(s string) string {
	var b strings.Builder
	for _, w := range strings.Fields(s) {
		b.WriteString(strings.ToUpper(w[:1]) + w[1:])
	}
	return b.String()
}

// snake converts "ActionItems" to "action_items".
func snake(s string) string {
	var b strings.Builder
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				b.WriteByte('_')
			}
			b.WriteRune(r - 'A' + 'a')
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// singular naively singularizes an English plural. Irregular cases are
// handled via the rename config.
func singular(s string) string {
	if strings.HasSuffix(s, "ses") || strings.HasSuffix(s, "ss") {
		return s
	}
	return strings.TrimSuffix(s, "s")
}

// lowerCamel converts an exported field name to a parameter name:
// UUID → uuid, ActionItemUUID → actionItemUUID, Name → name.
func lowerCamel(s string) string {
	// Lowercase the leading run of uppercase letters, keeping the last one
	// uppercase when it starts a new word (ActionItemUUID → actionItemUUID).
	upper := 0
	for upper < len(s) && s[upper] >= 'A' && s[upper] <= 'Z' {
		upper++
	}
	if upper == 0 {
		return s
	}
	if upper == len(s) {
		return strings.ToLower(s)
	}
	if upper > 1 {
		upper--
	}
	return strings.ToLower(s[:upper]) + s[upper:]
}

// aliasRegistry assigns public names to gen types, detecting collisions.
type aliasRegistry struct {
	byPublic map[string]string // public name → gen name
	byGen    map[string]string // gen name → public name
	order    []string          // public names in registration order
	rename   map[string]string // gen name → forced public name (config)
	reserved map[string]bool   // names taken by generated (non-alias) types
}

func newAliasRegistry(rename map[string]string) *aliasRegistry {
	return &aliasRegistry{
		byPublic: map[string]string{},
		byGen:    map[string]string{},
		rename:   rename,
		reserved: map[string]bool{},
	}
}

// add registers genName under public (or its configured rename). Returns
// the final public name. Registering the same pair twice is a no-op;
// claiming a public name already held by a different gen type is an error.
func (r *aliasRegistry) add(public, genName string) (string, error) {
	if forced, ok := r.rename[genName]; ok {
		public = forced
	}
	if existing, ok := r.byGen[genName]; ok {
		if existing != public {
			return "", fmt.Errorf("gen type %s already exposed as %s (requested %s)", genName, existing, public)
		}
		return public, nil
	}
	if held, ok := r.byPublic[public]; ok && held != genName {
		return "", fmt.Errorf("public name %s already taken by gen type %s (requested for %s); resolve via rename.aliases in .genwaroom.yaml", public, held, genName)
	}
	if r.reserved[public] {
		return "", fmt.Errorf("public name %s already taken by a generated type (requested for %s); resolve via rename.aliases in .genwaroom.yaml", public, genName)
	}
	r.byPublic[public] = genName
	r.byGen[genName] = public
	r.order = append(r.order, public)
	return public, nil
}

// addPreferred tries candidates in order and registers the first one that
// does not collide. A rename config entry for genName overrides them all.
func (r *aliasRegistry) addPreferred(genName string, candidates ...string) (string, error) {
	if forced, ok := r.rename[genName]; ok {
		return r.add(forced, genName)
	}
	var lastErr error
	for _, c := range candidates {
		public, err := r.add(c, genName)
		if err == nil {
			return public, nil
		}
		lastErr = err
	}
	return "", lastErr
}

// reserve claims a public name for a generated (non-alias) type such as
// XxxList or XxxOptions.
func (r *aliasRegistry) reserve(public string) error {
	if held, ok := r.byPublic[public]; ok {
		return fmt.Errorf("generated type name %s collides with alias of gen type %s", public, held)
	}
	if r.reserved[public] {
		return fmt.Errorf("generated type name %s is claimed twice", public)
	}
	r.reserved[public] = true
	return nil
}

// sortedAliases returns aliases sorted by public name for deterministic
// output.
func (r *aliasRegistry) sortedAliases() []Alias {
	publics := append([]string(nil), r.order...)
	sort.Strings(publics)
	aliases := make([]Alias, 0, len(publics))
	for _, p := range publics {
		aliases = append(aliases, Alias{Public: p, Gen: r.byPublic[p]})
	}
	return aliases
}
