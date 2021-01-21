package token

import "testing"

func TestBuiltins(t *testing.T) {
	if k := LookupIdent("a"); k != IDENT {
		t.Fatalf("CASE 'a': Expected Ident, got %s", k)
	}
	if k := LookupIdent("fn"); k != FUNCTION {
		t.Fatalf("CASE 'fn': Expected FUNCTION, got %s", k)
	}
}
