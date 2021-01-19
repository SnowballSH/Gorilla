package config

import "testing"

func TestBuiltins(t *testing.T) {
	if GetOSNewline("WINDOWS") != OSNEWLINES["windows"] {
		t.Fatalf("windows not match")
	}
}
