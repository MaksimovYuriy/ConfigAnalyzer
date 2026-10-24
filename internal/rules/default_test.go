package rules_test

import (
	"strings"
	"testing"

	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/analyzer"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/entity"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/rules"
	"github.com/MaksimovYuriy/ConfigAnalyzer/pkg/parser"
)

func TestDefaultRulesCanBeReplacedLocally(t *testing.T) {
	doc, err := parser.Parse(strings.NewReader("log: {level: debug}"))
	if err != nil {
		t.Fatal(err)
	}
	custom := rules.Default()
	custom.Fields[0] = func(entity.Field) []entity.Problem { return nil }
	if got := analyzer.Analyze(doc.Root, nil, custom); len(got) != 0 {
		t.Fatalf("local replacement did not disable debug check: %+v", got)
	}
	if got := analyzer.Analyze(doc.Root, nil, rules.Default()); len(got) != 1 || got[0].RuleID != "debug-logging" {
		t.Fatalf("local replacement changed default registry: %+v", got)
	}
}
