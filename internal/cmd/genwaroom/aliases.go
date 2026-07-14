package main

import (
	"fmt"
	"sort"
)

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
