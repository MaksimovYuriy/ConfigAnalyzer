package analyzer_test

import (
	"strings"
	"testing"

	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/analyzer"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/entity"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/rules"
	"github.com/MaksimovYuriy/ConfigAnalyzer/pkg/parser"
)

func TestAdditionalFileRuleRunsOnce(t *testing.T) {
	document, err := parser.Parse(strings.NewReader("logging: {level: debug}\nservices: [{name: one}, {name: two}]"))
	if err != nil {
		t.Fatal(err)
	}
	calls := 0
	custom := func(file entity.FileMetadata) []entity.Problem {
		calls++
		return []entity.Problem{{RuleID: "custom-file-rule", Path: file.Path}}
	}
	checks := rules.Default()
	checks.Files = append(checks.Files, custom)
	metadata := &entity.FileMetadata{Path: "config.yaml", Mode: 0600}
	problems := analyzer.Analyze(document.Root, metadata, checks)
	if calls != 1 || len(problems) != 2 || problems[0].RuleID != "custom-file-rule" || problems[1].RuleID != "debug-logging" {
		t.Fatalf("calls=%d, problems=%+v", calls, problems)
	}
	calls = 0
	metadata = nil
	problems = analyzer.Analyze(document.Root, metadata, checks)
	if calls != 0 || len(problems) != 1 || problems[0].RuleID != "debug-logging" {
		t.Fatalf("no-metadata analysis: calls=%d, problems=%+v", calls, problems)
	}
}

func TestFileAndFieldFindingsAreCombined(t *testing.T) {
	document, err := parser.Parse(strings.NewReader("logging: {level: debug}"))
	if err != nil {
		t.Fatal(err)
	}
	problems := analyzer.Analyze(document.Root, &entity.FileMetadata{Path: "config.yaml", Mode: 0666}, rules.Default())
	if len(problems) != 2 || problems[0].RuleID != "world-writable-config" || problems[1].RuleID != "debug-logging" {
		t.Fatalf("unexpected combined findings: %+v", problems)
	}
}
