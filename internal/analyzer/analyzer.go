package analyzer

import "go.yaml.in/yaml/v3"

func Analyze(root *yaml.Node, rules ...Rule) []Problem {
	var problems []Problem
	visit := func(field Field) {
		for _, rule := range rules {
			problems = append(problems, rule(field)...)
		}
	}
	walk(Field{Node: root, Path: "$"}, visit, make(map[*yaml.Node]bool))
	return problems
}
