package cli

import "testing"

func TestParseArgs(t *testing.T) {
	for _, tt := range []struct {
		name string
		args []string
		want options
	}{
		{"file then flag", []string{"config.yaml", "-s"}, options{path: "config.yaml", silent: true}},
		{"flag then file", []string{"--silent", "config.json"}, options{path: "config.json", silent: true}},
		{"stdin", []string{"--stdin", "-s"}, options{fromStdin: true, silent: true}},
		{"option-like filename", []string{"--", "--stdin"}, options{path: "--stdin"}},
		{"help", []string{"--help"}, options{help: true}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseArgs(tt.args)
			if err != nil || got != tt.want {
				t.Fatalf("parseArgs = %+v, %v; want %+v", got, err, tt.want)
			}
		})
	}
}

func TestParseArgsRejectsAmbiguousInput(t *testing.T) {
	for _, args := range [][]string{
		nil, {""}, {"", "config.yaml"}, {"--stdin", ""},
		{"a.yaml", "b.yaml"}, {"--stdin", "config.yaml"}, {"--unknown"},
	} {
		if _, err := parseArgs(args); err == nil {
			t.Errorf("expected error for %q", args)
		}
	}
}
