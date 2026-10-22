package parser_test

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/MaksimovYuriy/ConfigAnalyzer/pkg/parser"
	"go.yaml.in/yaml/v3"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		format parser.Format
		kind   yaml.Kind
		tag    string
		value  string
	}{
		{"JSON object", `{"log":{"level":"debug"}}`, parser.JSON, yaml.MappingNode, "!!map", ""},
		{"YAML object", "# comment\nlog:\n  level: debug\n", parser.YAML, yaml.MappingNode, "!!map", ""},
		{"JSON array", `[false, {"password":"secret"}]`, parser.JSON, yaml.SequenceNode, "!!seq", ""},
		{"YAML array", "- false\n- password: secret", parser.YAML, yaml.SequenceNode, "!!seq", ""},
		{"large JSON number", "1e400", parser.JSON, yaml.ScalarNode, "!!str", "1e400"},
		{"number root", "5", parser.JSON, yaml.ScalarNode, "!!int", "5"},
		{"YAML number root", "---\n5", parser.YAML, yaml.ScalarNode, "!!int", "5"},
		{"string root", "a configuration string", parser.YAML, yaml.ScalarNode, "!!str", "a configuration string"},
		{"quoted string root", `"hello"`, parser.JSON, yaml.ScalarNode, "!!str", "hello"},
		{"boolean root", "false", parser.JSON, yaml.ScalarNode, "!!bool", "false"},
		{"null root", "null", parser.JSON, yaml.ScalarNode, "!!null", "null"},
		{"explicit empty document", "---\n", parser.YAML, yaml.ScalarNode, "!!null", ""},
		{"YAML flow mapping", "{level: debug}", parser.YAML, yaml.MappingNode, "!!map", ""},
		{"leading blank lines", "\n\n# comment\nlevel: debug", parser.YAML, yaml.MappingNode, "!!map", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := parser.Parse(strings.NewReader(tt.input))
			if err != nil {
				t.Fatal(err)
			}
			if doc.Format != tt.format {
				t.Errorf("format = %q, want %q", doc.Format, tt.format)
			}
			if doc.Root.Kind != tt.kind || doc.Root.ShortTag() != tt.tag || doc.Root.Value != tt.value {
				t.Errorf("root = kind %v, tag %q, value %q; want %v, %q, %q", doc.Root.Kind, doc.Root.ShortTag(), doc.Root.Value, tt.kind, tt.tag, tt.value)
			}
		})
	}
}

func TestParsePreservesNestedValuesAndTags(t *testing.T) {
	input := "services:\n  - tls:\n      verify: false\n    password: !vault |\n      $ANSIBLE_VAULT;1.1;AES256\n      encrypted\n"
	doc, err := parser.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	service := doc.Root.Content[1].Content[0]
	verify := service.Content[1].Content[1]
	if verify.Tag != "!!bool" || verify.Value != "false" {
		t.Fatalf("TLS value = %#v", verify)
	}
	password := service.Content[3]
	if password.Tag != "!vault" || password.Line != 4 || !strings.Contains(password.Value, "$ANSIBLE_VAULT") {
		t.Fatalf("vault node lost metadata: %#v", password)
	}
}

func TestParsePreservesAliases(t *testing.T) {
	doc, err := parser.Parse(strings.NewReader("defaults: &defaults\n  verify: false\nservice: *defaults\n"))
	if err != nil {
		t.Fatal(err)
	}
	alias := doc.Root.Content[3]
	if alias.Kind != yaml.AliasNode || alias.Alias != doc.Root.Content[1] {
		t.Fatal("alias target was not preserved")
	}
}

func TestParseRejectsInvalidDocuments(t *testing.T) {
	tests := map[string]string{
		"empty":                    "",
		"whitespace":               " \n\t",
		"comments only":            "# comment\n",
		"unclosed sequence":        "values: [1, 2",
		"invalid indentation":      "key:\n  child: true\n other: false\n",
		"duplicate YAML key":       "password: first\npassword: second\n",
		"duplicate JSON key":       `{"password":"first","password":"second"}`,
		"multiple YAML documents":  "a: 1\n---\nb: 2\n",
		"trailing document marker": "a: 1\n---\n",
		"invalid second document":  "a: 1\n---\nb: [",
		"recursive alias":          "value: &loop [*loop]",
		"unknown alias":            "value: *missing",
	}
	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			doc, err := parser.Parse(strings.NewReader(input))
			if err == nil || doc != nil {
				t.Fatalf("Parse() = %v, %v; want nil document and error", doc, err)
			}
		})
	}
}

type failedReader struct{ err error }

func (r failedReader) Read([]byte) (int, error) { return 0, r.err }

func TestParseReadError(t *testing.T) {
	want := errors.New("input unavailable")
	_, err := parser.Parse(failedReader{err: want})
	if !errors.Is(err, want) {
		t.Fatalf("error = %v, want wrapped read error", err)
	}
}

func TestParseDoesNotCloseReader(t *testing.T) {
	r := &trackingReader{Reader: strings.NewReader("5")}
	if _, err := parser.Parse(r); err != nil {
		t.Fatal(err)
	}
	if r.closed {
		t.Fatal("parser closed caller's reader")
	}
}

type trackingReader struct {
	io.Reader
	closed bool
}

func (r *trackingReader) Close() error {
	r.closed = true
	return nil
}
