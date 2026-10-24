package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/cli"
)

func TestCLIFilePermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("requires Unix permission bits")
	}
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte("logging: {level: info}"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct {
		name    string
		mode    os.FileMode
		silent  bool
		code    int
		finding bool
	}{
		{"safe", 0644, false, cli.ExitOK, false},
		{"group writable", 0660, false, cli.ExitOK, false},
		{"public writable", 0666, false, cli.ExitFindings, true},
		{"silent", 0666, true, cli.ExitOK, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if err := os.Chmod(path, tt.mode); err != nil {
				t.Fatal(err)
			}
			args := []string{path}
			if tt.silent {
				args = append(args, "--silent")
			}
			var stdout, stderr bytes.Buffer
			code := cli.Run(args, strings.NewReader(""), &stdout, &stderr)
			if code != tt.code || stderr.Len() != 0 {
				t.Fatalf("code=%d, stderr=%s", code, &stderr)
			}
			found := strings.Contains(stdout.String(), "HIGH: "+path+":")
			if found != tt.finding {
				t.Fatalf("unexpected output: %s", &stdout)
			}
		})
	}
}
