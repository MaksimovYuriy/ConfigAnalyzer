package analyzer

import (
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/entity"
	"go.yaml.in/yaml/v3"
)

func Analyze(root *yaml.Node, rules ...entity.Rule) []entity.Problem {
	var problems []entity.Problem
	visit := func(field entity.Field) {
		for _, rule := range rules {
			problems = append(problems, rule(field)...)
		}
	}
	walk(entity.Field{Node: root, Path: "$"}, visit, make(map[*yaml.Node]bool))
	return problems
}
