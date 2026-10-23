package analyzer

import (
	"fmt"
	"strconv"
	"strings"

	"go.yaml.in/yaml/v3"
)

func walk(field Field, visit func(Field), active map[*yaml.Node]bool) {
	node := field.Node
	if node == nil || active[node] {
		return
	}
	active[node] = true
	defer delete(active, node)

	switch node.Kind {
	case yaml.AliasNode:
		field.Node = node.Alias
		walk(field, visit, active)
		return
	case yaml.DocumentNode:
		for _, child := range node.Content {
			field.Node = child
			walk(field, visit, active)
		}
		return
	}

	visit(field)

	switch node.Kind {
	case yaml.MappingNode:
		for i := 0; i < len(node.Content); i += 2 {
			key, value := node.Content[i], node.Content[i+1]
			if key.Kind != yaml.ScalarNode {
				continue
			}
			child := childField(field, value, key.Value, fieldPath(field.Path, key.Value))
			walk(child, visit, active)
		}
	case yaml.SequenceNode:
		for i, value := range node.Content {
			path := fmt.Sprintf("%s[%d]", field.Path, i)
			child := childField(field, value, "", path)
			walk(child, visit, active)
		}
	}
}

func childField(parent Field, node *yaml.Node, key, path string) Field {
	ancestors := parent.Ancestors
	if parent.Key != "" {
		ancestors = make([]string, len(parent.Ancestors)+1)
		copy(ancestors, parent.Ancestors)
		ancestors[len(parent.Ancestors)] = parent.Key
	}
	return Field{Node: node, Key: key, Path: path, Ancestors: ancestors}
}

func fieldPath(parent, key string) string {
	if key == "" || strings.ContainsAny(key, ".[]\"\\ \t\r\n") {
		return parent + "[" + strconv.Quote(key) + "]"
	}
	if parent == "$" {
		return key
	}
	return parent + "." + key
}
