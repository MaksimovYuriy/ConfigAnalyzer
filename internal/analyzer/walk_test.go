package analyzer_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/analyzer"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/entity"
	"github.com/MaksimovYuriy/ConfigAnalyzer/pkg/parser"
	"go.yaml.in/yaml/v3"
)

func TestWalkOrderAndContext(t *testing.T) {
	doc, err := parser.Parse(strings.NewReader("services:\n- log: {level: debug}\n- log: {level: info}\nlast: true"))
	if err != nil {
		t.Fatal(err)
	}
	var paths []string
	analyzer.Analyze(doc.Root, func(field entity.Field) []entity.Problem {
		paths = append(paths, field.Path)
		if field.Key == "level" && !reflect.DeepEqual(field.Ancestors, []string{"services", "log"}) {
			t.Errorf("unexpected ancestors: %v", field.Ancestors)
		}
		return nil
	})
	want := []string{"$", "services", "services[0]", "services[0].log", "services[0].log.level", "services[1]", "services[1].log", "services[1].log.level", "last"}
	if !reflect.DeepEqual(paths, want) {
		t.Fatalf("paths = %v, want %v", paths, want)
	}
}

func TestWalkDeepTree(t *testing.T) {
	const depth = 2048
	root := &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: "leaf"}
	for i := 0; i < depth; i++ {
		root = &yaml.Node{Kind: yaml.SequenceNode, Content: []*yaml.Node{root}}
	}
	count := 0
	analyzer.Analyze(root, func(field entity.Field) []entity.Problem {
		count++
		if field.Node.Kind == yaml.ScalarNode && field.Path != "$"+strings.Repeat("[0]", depth) {
			t.Error("incorrect leaf path")
		}
		return nil
	})
	if count != depth+1 {
		t.Fatalf("visited %d nodes, want %d", count, depth+1)
	}
}

func TestWalkCyclicAlias(t *testing.T) {
	root := &yaml.Node{Kind: yaml.SequenceNode}
	root.Content = []*yaml.Node{{Kind: yaml.AliasNode, Alias: root}}
	count := 0
	analyzer.Analyze(root, func(entity.Field) []entity.Problem {
		count++
		return nil
	})
	if count != 1 {
		t.Fatalf("visited %d values in cyclic tree", count)
	}
}

func TestWalkNilRoot(t *testing.T) {
	analyzer.Analyze(nil, func(entity.Field) []entity.Problem {
		t.Fatal("rule called for nil root")
		return nil
	})
}
