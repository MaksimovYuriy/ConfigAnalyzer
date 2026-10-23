package helpers

import (
	"strings"
	"unicode"
)

type Names map[string]struct{}

func NewNames(aliases ...string) Names {
	names := make(Names, len(aliases))
	for _, alias := range aliases {
		if name := Normalize(alias); name != "" {
			names[name] = struct{}{}
		}
	}
	return names
}

func (names Names) Matches(name string) bool {
	_, ok := names[Normalize(name)]
	return ok
}

func Normalize(name string) string {
	return strings.Map(func(r rune) rune {
		if r == '-' || r == '_' {
			return -1
		}
		return unicode.ToLower(r)
	}, name)
}
