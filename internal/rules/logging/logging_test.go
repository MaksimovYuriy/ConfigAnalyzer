package logging_test

import (
	"strings"
	"testing"

	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/analyzer"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/rules/logging"
	"github.com/MaksimovYuriy/ConfigAnalyzer/pkg/parser"
)

func TestDebugLogging(t *testing.T) {
	tests := []struct {
		name, input string
		paths       []string
	}{
		{"JSON", `{"log":{"level":"debug"}}`, []string{"log.level"}},
		{"array", "services:\n- logging:\n    level: DEBUG\n- logging:\n    level: info", []string{"services[0].logging.level"}},
		{"aliases", "logging-level: debug\nlogLevel: debug", []string{"logging-level", "logLevel"}},
		{"nested context", "logger:\n  output:\n    level: debug", []string{"logger.output.level"}},
		{"alias target", "log: &log\n  level: debug\nservice:\n  logging: *log", []string{"log.level", "service.logging.level"}},
		{"escaped path", "'a.b':\n  log_level: debug", []string{`$["a.b"].log_level`}},
		{"unrelated level", "game:\n  level: debug", nil},
		{"safe", "logging:\n  level: info", nil},
		{"boolean", "log:\n  level: true", nil},
		{"scalar", "5", nil},
		{"tagged secret", "logging:\n  level: !vault debug", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := parser.Parse(strings.NewReader(tt.input))
			if err != nil {
				t.Fatal(err)
			}
			problems := analyzer.Analyze(doc.Root, logging.Check)
			if len(problems) != len(tt.paths) {
				t.Fatalf("problems = %+v, want paths %v", problems, tt.paths)
			}
			for i, problem := range problems {
				if problem.Path != tt.paths[i] || problem.Severity != analyzer.Low || problem.RuleID != "debug-logging" || problem.Recommendation == "" || problem.Explanation == "" {
					t.Errorf("unexpected problem: %+v", problem)
				}
			}
		})
	}
}
