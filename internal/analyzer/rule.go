package analyzer

import "go.yaml.in/yaml/v3"

type Field struct {
	Key       string
	Path      string
	Ancestors []string
	Node      *yaml.Node
}

type Rule func(Field) []Problem
