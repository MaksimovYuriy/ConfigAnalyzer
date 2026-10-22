package rules

import "testing"

func TestNameSetMatches(t *testing.T) {
	names := newNameSet("tls_enabled", "enable_tls")
	for _, name := range []string{"tls_enabled", "tls-enabled", "tlsEnabled", "TLSEnabled", "enableTLS"} {
		if !names.matches(name) {
			t.Errorf("expected match for %q", name)
		}
	}
	for _, name := range []string{"", "tls", "my_tls_enabled", "tls.enabled", "tls enabled", "tls_enabled_extra"} {
		if names.matches(name) {
			t.Errorf("unexpected match for %q", name)
		}
	}
}

func TestNameSetDeduplicatesAliases(t *testing.T) {
	names := newNameSet("password", "PASSWORD", "pass-word", "", "_-")
	if len(names) != 1 || !names.matches("password") {
		t.Fatalf("expected one normalized nonempty alias, got %v", names)
	}
}

func TestNameSetZeroValue(t *testing.T) {
	var names nameSet
	if names.matches("password") {
		t.Fatal("empty name set matched a parameter")
	}
}

func TestParameterAliases(t *testing.T) {
	tests := []struct {
		name  string
		names nameSet
		alias string
	}{
		{"password synonym", passwordNames, "passwd"},
		{"password spelling", passwordNames, "dbPassword"},
		{"logging context", loggingNames, "LOGGER"},
		{"log level", logLevelNames, "logLevel"},
		{"TLS enabled", tlsEnabledNames, "tlsEnabled"},
		{"TLS verification", tlsVerifyNames, "verifySSL"},
		{"TLS verification bypass", tlsSkipVerifyNames, "insecureSkipVerify"},
		{"listen address", listenAddressNames, "bind-address"},
		{"algorithm", algorithmNames, "digest-algorithm"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.names.matches(tt.alias) {
				t.Errorf("expected match for %q", tt.alias)
			}
		})
	}
	if tlsVerifyNames.matches("insecureSkipVerify") || tlsEnabledNames.matches("verifySSL") {
		t.Fatal("TLS parameters with different meanings must remain separate")
	}
	if passwordNames.matches("password_policy") || algorithmNames.matches("description") {
		t.Fatal("unrelated parameters must not match")
	}
}
