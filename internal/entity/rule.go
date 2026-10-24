package entity

import "go.yaml.in/yaml/v3"

type Field struct {
	Key       string
	Path      string
	Ancestors []string
	Node      *yaml.Node
}

type FieldRule func(Field) []Problem

type FileRule func(FileMetadata) []Problem

type RuleSet struct {
	Fields []FieldRule
	Files  []FileRule
}
