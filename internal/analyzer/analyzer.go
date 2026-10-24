package analyzer

import (
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/entity"
	"go.yaml.in/yaml/v3"
)

func Analyze(root *yaml.Node, fileMetadata *entity.FileMetadata, rules entity.RuleSet) []entity.Problem {
	var problems []entity.Problem
	if fileMetadata != nil {
		for _, rule := range rules.Files {
			problems = append(problems, rule(*fileMetadata)...)
		}
	}
	checkField := func(field entity.Field) {
		for _, rule := range rules.Fields {
			problems = append(problems, rule(field)...)
		}
	}
	walk(entity.Field{Node: root, Path: "$"}, checkField, make(map[*yaml.Node]bool))
	return problems
}
