package algorithm

import (
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/entity"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/rules/helpers"
)

var (
	insecureAlgorithms = helpers.NewNames(
		"MD2",
		"MD4",
		"MD5",
		"SHA1",
		"DES",
		"3DES",
		"TripleDES",
		"DES-EDE3",
		"DES-EDE3-CBC",
		"RC4",
	)
	algorithmNames = helpers.NewNames(
		"algorithm",
		"digest_algorithm",
		"hash_algorithm",
		"hashing_algorithm",
		"cipher",
		"algorithms",
		"ciphers",
		"cipher_algorithm",
		"encryption_algorithm",
		"signature_algorithm",
		"hash",
	)
)

func Check(field entity.Field) []entity.Problem {
	if !helpers.NamedValue(field, algorithmNames) {
		return nil
	}
	value, isString := helpers.ScalarText(field)
	if !isString {
		return nil
	}
	if !insecureAlgorithms.Matches(value) {
		return nil
	}
	return []entity.Problem{{
		RuleID:         "insecure-algorithm",
		Path:           field.Path,
		Severity:       entity.High,
		Explanation:    "Выбран устаревший или небезопасный криптографический алгоритм.",
		Recommendation: "Замените алгоритм на современный, подходящий для этой задачи.",
	}}
}
