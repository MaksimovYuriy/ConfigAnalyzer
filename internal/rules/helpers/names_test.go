package helpers

import "testing"

func TestNameSetMatches(t *testing.T) {
	names := NewNames("tls_enabled", "enable_tls")
	for _, name := range []string{"tls_enabled", "tls-enabled", "tlsEnabled", "TLSEnabled", "enableTLS"} {
		if !names.Matches(name) {
			t.Errorf("expected match for %q", name)
		}
	}
	for _, name := range []string{"", "tls", "my_tls_enabled", "tls.enabled", "tls enabled", "tls_enabled_extra"} {
		if names.Matches(name) {
			t.Errorf("unexpected match for %q", name)
		}
	}
}

func TestNameSetDeduplicatesAliases(t *testing.T) {
	names := NewNames("password", "PASSWORD", "pass-word", "", "_-")
	if len(names) != 1 || !names.Matches("password") {
		t.Fatalf("expected one normalized nonempty alias, got %v", names)
	}
}

func TestNameSetZeroValue(t *testing.T) {
	var names Names
	if names.Matches("password") {
		t.Fatal("empty name set matched a parameter")
	}
}
