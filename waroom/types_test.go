package waroom

import "testing"

func TestOptNilIntConstructors(t *testing.T) {
    correct := NewOptNilInt(1)
    legacy := NewOptNillInt(1)
    if !correct.Set || !legacy.Set {
        t.Fatal("both constructors must produce a set value")
    }
    if correct.Value != legacy.Value {
        t.Fatalf("constructors disagree: %d != %d", correct.Value, legacy.Value)
    }
}
