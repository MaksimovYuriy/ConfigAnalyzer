package password_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/analyzer"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/entity"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/rules/password"
	"github.com/MaksimovYuriy/ConfigAnalyzer/pkg/parser"
)

func TestCheck(t *testing.T) {
	for _, tt := range []struct {
		name, input, id, path string
		severity              entity.Severity
	}{
		{"password", "db: {password: super-secret}", "plaintext-password", "db.password", entity.High},
		{"password alias", "dbPassword: 12345", "plaintext-password", "dbPassword", entity.High},
		{"password list", "password: [secret]", "plaintext-password", "password[0]", entity.High},
		{"partial reference", "password: prefix-${SECRET}", "plaintext-password", "password", entity.High},
		{"reference with literal suffix", "password: '{{ secret }}suffix'", "plaintext-password", "password", entity.High},
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
		"password: '${DB_PASSWORD}'",
		"password: '{{ vault_db_password }}'",
		"password: !vault |\n  $ANSIBLE_VAULT;1.1;AES256\n  abcdef\n",
		"password: !env DB_PASSWORD",
		"password: !secret db_password",
		"password: ''",
		"password: null",
		"password_policy: strong",
		"5",
	} {
		t.Run(input, func(t *testing.T) {
			if problems := audit(t, input); len(problems) != 0 {
				t.Fatalf("unexpected findings: %+v", problems)
			}
		})
	}
}

func TestSecurityRulesInAliases(t *testing.T) {
	problems := audit(t, "defaults: &settings\n  password: secret\nservice: *settings\n")
	if len(problems) != 2 || problems[1].Path != "service.password" {
		t.Fatalf("alias findings: %+v", problems)
	}
}

func audit(t *testing.T, input string) []entity.Problem {
	t.Helper()
	doc, err := parser.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	return analyzer.Analyze(doc.Root, nil, entity.RuleSet{Fields: []entity.FieldRule{password.Check}})
}

func TestSecretReferences(t *testing.T) {
	for _, tt := range []struct {
		value     string
		reference bool
	}{
		{"${DB_PASSWORD}", true},
		{"{{ vault_db_password }}", true},
		{"{{ lookup('env', 'DB_PASSWORD') }}", true},
		{"${}", false},
		{"${   }", false},
		{"{{}}", false},
		{"{{   }}", false},
		{"${", false},
		{"{{}", false},
		{"prefix-${DB_PASSWORD}", false},
		{"${DB_PASSWORD}-suffix", false},
		{"${A}${B}", false},
		{"{{ secret }}suffix", false},
		{"${A\nB}", false},
	} {
		t.Run(tt.value, func(t *testing.T) {
			problems := audit(t, "password: "+strconv.Quote(tt.value))
			if (len(problems) == 0) != tt.reference {
				t.Fatalf("reference=%v, findings=%+v", tt.reference, problems)
			}
		})
	}
}
