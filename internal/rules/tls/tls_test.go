package tls_test

import (
	"strings"
	"testing"

	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/analyzer"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/rules/tls"
	"github.com/MaksimovYuriy/ConfigAnalyzer/pkg/parser"
)

func TestCheck(t *testing.T) {
	for _, tt := range []struct {
		name, input, id, path string
		severity              analyzer.Severity
	}{
		{"TLS scalar", "tls: false", "tls-disabled", "tls", analyzer.High},
		{"TLS enabled", "tlsEnabled: false", "tls-disabled", "tlsEnabled", analyzer.High},
		{"nested TLS", "server: {tls: {enabled: false}}", "tls-disabled", "server.tls.enabled", analyzer.High},
		{"disable TLS", "disable-tls: true", "tls-disabled", "disable-tls", analyzer.High},
		{"verify SSL", "verify_ssl: false", "tls-verification-disabled", "verify_ssl", analyzer.High},
		{"nested verify", "tls: {verify: false}", "tls-verification-disabled", "tls.verify", analyzer.High},
		{"skip verify", "insecureSkipVerify: true", "tls-verification-disabled", "insecureSkipVerify", analyzer.High},
		{"quoted boolean", "ssl_verify: 'off'", "tls-verification-disabled", "ssl_verify", analyzer.High},
		{"numeric boolean", "tls_enabled: 0", "tls-disabled", "tls_enabled", analyzer.High},
		{"reject unauthorized", "rejectUnauthorized: false", "tls-verification-disabled", "rejectUnauthorized", analyzer.High},
	} {
		t.Run(tt.name, func(t *testing.T) {
			problems := audit(t, tt.input)
			if len(problems) != 1 {
				t.Fatalf("problems = %+v, want one", problems)
			}
			p := problems[0]
			if p.RuleID != tt.id || p.Path != tt.path || p.Severity != tt.severity || p.Explanation == "" || p.Recommendation == "" {
				t.Fatalf("unexpected finding: %+v", p)
			}
		})
	}
}

func TestSafeOrUnrelatedValues(t *testing.T) {
	for _, input := range []string{
		"tls: true",
		"tls: {enabled: true, verify: true}",
		"tls_enabled: '${TLS_ENABLED}'",
		"verify_ssl: true",
		"insecure_skip_verify: false",
		"disable_tls: false",
		"feature: {enabled: false, verify: false}",
		"tls: {enabled: !vault false}",
		"rejectUnauthorized: true",
		"5",
	} {
		t.Run(input, func(t *testing.T) {
			if problems := audit(t, input); len(problems) != 0 {
				t.Fatalf("unexpected findings: %+v", problems)
			}
		})
	}
}

func audit(t *testing.T, input string) []analyzer.Problem {
	t.Helper()
	doc, err := parser.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	return analyzer.Analyze(doc.Root, tls.Check)
}
