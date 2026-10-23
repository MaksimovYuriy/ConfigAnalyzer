package listen_test

import (
	"strings"
	"testing"

	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/analyzer"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/rules/listen"
	"github.com/MaksimovYuriy/ConfigAnalyzer/pkg/parser"
)

func TestCheck(t *testing.T) {
	for _, tt := range []struct {
		name, input, id, path string
		severity              analyzer.Severity
	}{
		{"wildcard", "host: 0.0.0.0", "wildcard-listen", "host", analyzer.Medium},
		{"wildcard with port", "listen_address: '0.0.0.0:8080'", "wildcard-listen", "listen_address", analyzer.Medium},
		{"listen list", "listen: ['127.0.0.1:80', '0.0.0.0:80']", "wildcard-listen", "listen[1]", analyzer.Medium},
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
		"host: 127.0.0.1",
		"listen: '127.0.0.1:8080'",
		"description: 0.0.0.0",
		"host: 'http://0.0.0.0:80'",
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
	return analyzer.Analyze(doc.Root, listen.Check)
}
