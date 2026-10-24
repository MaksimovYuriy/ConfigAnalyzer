package algorithm_test

import (
	"strings"
	"testing"

	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/analyzer"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/entity"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/rules/algorithm"
	"github.com/MaksimovYuriy/ConfigAnalyzer/pkg/parser"
)

func TestCheck(t *testing.T) {
	for _, tt := range []struct {
		name, input, id, path string
		severity              entity.Severity
	}{
		{"MD5", "storage: {digest-algorithm: MD5}", "insecure-algorithm", "storage.digest-algorithm", entity.High},
		{"SHA1 normalization", "hash_algorithm: sha-1", "insecure-algorithm", "hash_algorithm", entity.High},
		{"algorithm array", "ciphers: [AES-GCM, RC4]", "insecure-algorithm", "ciphers[1]", entity.High},
	} {
		t.Run(tt.name, func(t *testing.T) {
			problems := audit(t, tt.input)
			if len(problems) != 1 {
				t.Fatalf("problems = %+v, want one", problems)
			}
			p := problems[0]
			if p.RuleID != tt.id || p.Path != tt.path || p.Severity != tt.severity || p.Explanation == "" || p.Recommendation == "" {
				t.Fatalf("unexpected finding: %+v", p)
			}
		})
	}
}

func TestSafeOrUnrelatedValues(t *testing.T) {
	for _, input := range []string{
		"algorithm: SHA256",
		"algorithm: HMAC-SHA1",
		"description: MD5",
		"algorithm: 'MD5 is deprecated'",
		"cipher: AES-GCM",
		"5",
	} {
		t.Run(input, func(t *testing.T) {
			if problems := audit(t, input); len(problems) != 0 {
				t.Fatalf("unexpected findings: %+v", problems)
			}
		})
	}
}

func TestInsecureAlgorithmList(t *testing.T) {
	for _, name := range []string{"MD2", "MD4", "MD5", "SHA-1", "DES", "3DES", "TripleDES", "DES-EDE3", "DES-EDE3-CBC", "RC4"} {
		if problems := audit(t, "algorithm: "+name); len(problems) != 1 || problems[0].RuleID != "insecure-algorithm" {
			t.Errorf("algorithm %q: %+v", name, problems)
		}
	}
}

func audit(t *testing.T, input string) []entity.Problem {
	t.Helper()
	doc, err := parser.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	return analyzer.Analyze(doc.Root, algorithm.Check)
}
