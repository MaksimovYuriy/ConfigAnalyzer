package analyzer_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/analyzer"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/entity"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/rules"
	"github.com/MaksimovYuriy/ConfigAnalyzer/pkg/parser"
)

func TestAnalyzeAppliesIndependentRules(t *testing.T) {
	doc, err := parser.Parse(strings.NewReader("key: value"))
	if err != nil {
		t.Fatal(err)
	}
	check := func(id string) entity.FieldRule {
		return func(field entity.Field) []entity.Problem {
			if field.Key != "key" {
				return nil
			}
			return []entity.Problem{{RuleID: id, Path: field.Path}}
		}
	}
	problems := analyzer.Analyze(doc.Root, nil, entity.RuleSet{Fields: []entity.FieldRule{check("first"), check("second")}})
	if len(problems) != 2 || problems[0].RuleID != "first" || problems[1].RuleID != "second" {
		t.Fatalf("independent rule results: %+v", problems)
	}
}

func TestAnalyzeWithAdditionalRule(t *testing.T) {
	for _, input := range []string{
		"logging: {level: debug}\nservices:\n- session:\n    cookie:\n      http_only: false\n- session:\n    cookie:\n      http_only: true\n",
		`{"logging":{"level":"debug"},"services":[{"session":{"cookie":{"http_only":false}}},{"session":{"cookie":{"http_only":true}}}]}`,
	} {
		doc, err := parser.Parse(strings.NewReader(input))
		if err != nil {
			t.Fatal(err)
		}
		cookieRule := func(field entity.Field) []entity.Problem {
			if field.Key != "http_only" || field.Node.Tag != "!!bool" || field.Node.Value != "false" ||
				len(field.Ancestors) == 0 || field.Ancestors[len(field.Ancestors)-1] != "cookie" {
				return nil
			}
			return []entity.Problem{{
				RuleID: "test-cookie-http-only", Path: field.Path, Severity: entity.Medium,
				Explanation: "Cookie доступна JavaScript.", Recommendation: "Включите HttpOnly.",
			}}
		}
		builtins := rules.Default()
		before := analyzer.Analyze(doc.Root, nil, builtins)
		extended := rules.Default()
		extended.Fields = append(extended.Fields, cookieRule)
		got := analyzer.Analyze(doc.Root, nil, extended)
		want := append(append([]entity.Problem(nil), before...), entity.Problem{
			RuleID: "test-cookie-http-only", Path: "services[0].session.cookie.http_only", Severity: entity.Medium,
			Explanation: "Cookie доступна JavaScript.", Recommendation: "Включите HttpOnly.",
		})
		if len(before) != 1 || !reflect.DeepEqual(got, want) {
			t.Fatalf("extended findings = %+v; want %+v", got, want)
		}
		if after := analyzer.Analyze(doc.Root, nil, builtins); !reflect.DeepEqual(after, before) {
			t.Fatalf("extension changed built-in checks: %+v", after)
		}
		reversedRules := rules.Default()
		reversedRules.Fields = append([]entity.FieldRule{cookieRule}, reversedRules.Fields...)
		reversed := analyzer.Analyze(doc.Root, nil, reversedRules)
		if !reflect.DeepEqual(reversed, got) {
			t.Fatalf("changing rule order changed findings: %+v", reversed)
		}
	}
}
