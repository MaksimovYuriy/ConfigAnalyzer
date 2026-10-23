package helpers

import (
	"strings"

	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/analyzer"
	"go.yaml.in/yaml/v3"
)

func InContext(field analyzer.Field, names Names) bool {
	for _, name := range field.Ancestors {
		if names.Matches(name) {
			return true
		}
	}
	return false
}

func NamedValue(field analyzer.Field, names Names) bool {
	if field.Key != "" {
		return names.Matches(field.Key)
	}
	return len(field.Ancestors) > 0 && names.Matches(field.Ancestors[len(field.Ancestors)-1])
}

func ScalarText(field analyzer.Field) (string, bool) {
	if field.Node.Kind != yaml.ScalarNode || field.Node.ShortTag() != "!!str" {
		return "", false
	}
	return strings.TrimSpace(field.Node.Value), true
}

func BooleanValue(field analyzer.Field) (bool, bool) {
	if field.Node.Kind != yaml.ScalarNode {
		return false, false
	}
	switch field.Node.ShortTag() {
	case "!!bool", "!!str", "!!int":
	default:
		return false, false
	}
	switch strings.ToLower(strings.TrimSpace(field.Node.Value)) {
	case "true", "yes", "on", "1":
		return true, true
	case "false", "no", "off", "0":
		return false, true
	default:
		return false, false
	}
}
