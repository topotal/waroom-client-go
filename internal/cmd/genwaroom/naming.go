package main

import (
	"fmt"
	"strings"
)

// Conventions shared by both stages: analyze uses them to classify and
// decide, render uses them for the fixed mechanical joins.

// errorStatus maps ogen variant name suffixes to HTTP status codes. This
// is a closed set: an unknown suffix aborts generation.
var errorStatus = map[string]int{
	"BadRequest":          400,
	"Unauthorized":        401,
	"PaymentRequired":     402,
	"NotFound":            404,
	"UnprocessableEntity": 422,
}

// paginationFields is the exact shape a list ResponseMetadata must have.
var paginationFields = []string{"CurrentPage", "NextPage", "PrevPage", "TotalPages", "TotalCount"}

// unwrapVerbs maps operation verbs to the participle used when naming an
// unwrapped response type: createServiceLabel + Label → CreatedLabel.
var unwrapVerbs = map[string]string{
	"Create": "Created",
	"Update": "Updated",
}

// statusOf resolves the HTTP status of an error variant from its name
// suffix.
func statusOf(variant, rawName string) (int, error) {
	suffix := strings.TrimPrefix(variant, rawName)
	status, ok := errorStatus[suffix]
	if !ok {
		return 0, fmt.Errorf("error variant %s has unknown suffix %q", variant, suffix)
	}
	return status, nil
}

// Derived service names (fixed joins on the tag).

func serviceName(tag string) string { return singular(pascal(tag)) + "Service" }
func serviceAccessor(tag string) string { return pascal(tag) }
func serviceFileName(tag string) string { return snake(pascal(tag)) + "_gen.go" }

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

// sentence trims a spec summary and terminates it with a period so gofmt
// does not reformat the doc comment line as a heading.
func sentence(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	if !strings.HasSuffix(s, ".") && !strings.HasSuffix(s, "!") && !strings.HasSuffix(s, "?") {
		s += "."
	}
	return s
}

func leadingWord(s string) string {
	for i := 1; i < len(s); i++ {
		if s[i] >= 'A' && s[i] <= 'Z' {
			return s[:i]
		}
	}
	return s
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
