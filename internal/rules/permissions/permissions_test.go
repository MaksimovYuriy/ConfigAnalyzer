package permissions_test

import (
	"io/fs"
	"testing"

	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/entity"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/rules/permissions"
)

func TestCheck(t *testing.T) {
	for _, tt := range []struct {
		name    string
		mode    fs.FileMode
		finding bool
	}{
		{"owner only", 0600, false},
		{"public read", 0644, false},
		{"group write", 0660, false},
		{"others write", 0602, true},
		{"public write", 0666, true},
		{"all permissions", 0777, true},
		{"directory", fs.ModeDir | 0777, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			problems := permissions.Check(entity.FileMetadata{Path: "config.yaml", Mode: tt.mode})
			if !tt.finding {
				if len(problems) != 0 {
					t.Fatalf("unexpected findings: %+v", problems)
				}
				return
			}
			if len(problems) != 1 {
				t.Fatalf("expected one finding, got %+v", problems)
			}
			problem := problems[0]
			if problem.RuleID != "world-writable-config" || problem.Path != "config.yaml" || problem.Severity != entity.High || problem.Explanation == "" || problem.Recommendation == "" {
				t.Fatalf("unexpected problem: %+v", problem)
			}
		})
	}
}
