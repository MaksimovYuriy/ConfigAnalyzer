package tls_test

import (
	"strings"
	"testing"

	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/analyzer"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/entity"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/rules/tls"
	"github.com/MaksimovYuriy/ConfigAnalyzer/pkg/parser"
)

func TestCheck(t *testing.T) {
	for _, tt := range []struct {
		name, input, id, path string
		severity              entity.Severity
	}{
		{"TLS scalar", "tls: false", "tls-disabled", "tls", entity.High},
		{"TLS enabled", "tlsEnabled: false", "tls-disabled", "tlsEnabled", entity.High},
		{"nested TLS", "server: {tls: {enabled: false}}", "tls-disabled", "server.tls.enabled", entity.High},
		{"disable TLS", "disable-tls: true", "tls-disabled", "disable-tls", entity.High},
		{"verify SSL", "verify_ssl: false", "tls-verification-disabled", "verify_ssl", entity.High},
		{"nested verify", "tls: {verify: false}", "tls-verification-disabled", "tls.verify", entity.High},
		{"skip verify", "insecureSkipVerify: true", "tls-verification-disabled", "insecureSkipVerify", entity.High},
		{"quoted boolean", "ssl_verify: 'off'", "tls-verification-disabled", "ssl_verify", entity.High},
		{"numeric boolean", "tls_enabled: 0", "tls-disabled", "tls_enabled", entity.High},
		{"reject unauthorized", "rejectUnauthorized: false", "tls-verification-disabled", "rejectUnauthorized", entity.High},
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

func audit(t *testing.T, input string) []entity.Problem {
	t.Helper()
	doc, err := parser.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	return analyzer.Analyze(doc.Root, nil, entity.RuleSet{Fields: []entity.FieldRule{tls.Check}})
}
