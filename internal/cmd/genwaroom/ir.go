package main

// This file defines the semantic intermediate representation — the schema
// of api/sdk-model.yaml and the contract between the two generator
// stages. It holds only what `render` cannot derive on its own: facts
// read from the ogen package (type names, shapes) and naming decisions
// that went through heuristics or collision resolution. Everything that
// is a fixed mechanical join (file names, argument casing, constructor
// names, error status from the variant suffix) is derived at render time
// from these fields.

// Model is the root of the semantic model.
type Model struct {
	Services []*Service `yaml:"services"`
	Types    *TypesFile `yaml:"types"`
}

// Service groups operations by their spec tag. Service type name,
// accessor name, and file name are derived from Tag.
type Service struct {
	Tag string `yaml:"tag"` // "Labels"
	Ops []*Op  `yaml:"operations"`
}

// Op is one generated public method (plus its optional Options struct and
// All iterator).
type Op struct {
	ID      string `yaml:"id"`               // operationId: createServiceLabel
	RawName string `yaml:"rawName"`          // gen client method name
	GoName  string `yaml:"goName,omitempty"` // public method name; empty = RawName
	Method  string `yaml:"method"`           // POST
	Path    string `yaml:"path"`             // /api/v0/services/{name}/labels
	Summary string `yaml:"summary,omitempty"`

	PathParams []PathParam `yaml:"pathParams,omitempty"` // in path order
	ParamsGen  string      `yaml:"paramsGen,omitempty"`  // gen params struct: CreateServiceLabelParams
	Request    *Request    `yaml:"request,omitempty"`
	Options    *Options    `yaml:"options,omitempty"`
	Success    Success     `yaml:"success"`
	// Errors lists the error variant type names; the HTTP status is
	// derived from the variant name suffix (see errorStatus).
	Errors []string `yaml:"errors"`
	All    bool     `yaml:"all,omitempty"` // generate the auto-pagination iterator
}

// PathParam is a path placeholder promoted to a positional method
// argument. The argument name is derived from GenField.
type PathParam struct {
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

// NewFunc re-exports an Opt/Nil constructor; its name is New<Type>.
type NewFunc struct {
	Param string `yaml:"param"` // public param type: int or ListIncidentsSortOrder
	Type  string `yaml:"type"`  // public result type: OptInt
}

// Enum re-exports the constants of an aliased enum type.
type Enum struct {
	Type   string   `yaml:"type"` // ListIncidentsSortOrder
	Consts []string `yaml:"consts"`
}
