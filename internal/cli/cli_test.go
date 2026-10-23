package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/cli"
)

func TestRun(t *testing.T) {
	for _, tt := range []struct {
		name   string
		args   []string
		input  string
		code   int
		output string
	}{
		{"finding", []string{"--stdin"}, "log: {level: debug}", cli.ExitFindings, "LOW: log.level:"},
		{"silent", []string{"--stdin", "--silent"}, "log: {level: debug}", cli.ExitOK, "LOW: log.level:"},
		{"short silent", []string{"-s", "--stdin"}, "log: {level: debug}", cli.ExitOK, "LOW: log.level:"},
		{"safe", []string{"--stdin"}, "log: {level: info}", cli.ExitOK, "Проблем не найдено"},
		{"scalar", []string{"--stdin"}, "5", cli.ExitOK, "Проблем не найдено"},
		{"parse error despite silent", []string{"--stdin", "-s"}, "log: [", cli.ExitError, ""},
		{"missing input", nil, "", cli.ExitError, ""},
		{"conflicting inputs", []string{"file.yaml", "--stdin"}, "", cli.ExitError, ""},
		{"multiple files", []string{"a.yaml", "b.yaml"}, "", cli.ExitError, ""},
		{"unknown flag", []string{"--unknown"}, "", cli.ExitError, ""},
		{"help", []string{"--help"}, "", cli.ExitOK, "Usage:"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := cli.Run(tt.args, strings.NewReader(tt.input), &stdout, &stderr)
			if code != tt.code {
				t.Fatalf("code = %d, want %d; stderr: %s", code, tt.code, &stderr)
			}
			if !strings.Contains(stdout.String(), tt.output) {
				t.Errorf("stdout = %q, want %q", stdout.String(), tt.output)
			}
			if (code == cli.ExitError) != (stderr.Len() > 0) {
				t.Errorf("unexpected stderr: %q", stderr.String())
			}
		})
	}
}

func TestRunFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte("logging: {level: debug}"), 0600); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := cli.Run([]string{path, "-s"}, strings.NewReader("ignored"), &stdout, &stderr)
	if code != cli.ExitOK || !strings.Contains(stdout.String(), "LOW: logging.level:") || stderr.Len() != 0 {
		t.Fatalf("code=%d stdout=%q stderr=%q", code, &stdout, &stderr)
	}
	code = cli.Run([]string{filepath.Join(t.TempDir(), "missing.yaml"), "-s"}, strings.NewReader(""), &stdout, &stderr)
	if code != cli.ExitError {
		t.Fatalf("missing file returned %d", code)
	}
}
