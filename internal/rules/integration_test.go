package rules_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/analyzer"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/rules"
	"github.com/MaksimovYuriy/ConfigAnalyzer/pkg/parser"
)

func TestAllRulesJSONAndNoPasswordLeak(t *testing.T) {
	problems := audit(t, `{"log":{"level":"debug"},"password":"do-not-print-this","tls_enabled":false,"verify_ssl":false,"host":"0.0.0.0","algorithm":"MD5"}`)
	if len(problems) != 6 {
		t.Fatalf("got %d problems: %+v", len(problems), problems)
	}
	encoded, err := json.Marshal(problems)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(encoded), "do-not-print-this") {
		t.Fatal("password leaked into findings")
	}
}

func audit(t *testing.T, input string) []analyzer.Problem {
	t.Helper()
	doc, err := parser.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	return analyzer.Analyze(doc.Root, rules.Default()...)
}
