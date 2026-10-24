package password

import (
	"strings"

	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/entity"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/rules/helpers"
	"go.yaml.in/yaml/v3"
)

var (
	passwordNames = helpers.NewNames(
		"password",
		"passwd",
		"pwd",
		"db_password",
		"database_password",
		"user_password",
		"admin_password",
		"db_pass",
		"passphrase",
	)
)

func Check(field entity.Field) []entity.Problem {
	if !helpers.NamedValue(field, passwordNames) || field.Node.Kind != yaml.ScalarNode {
		return nil
	}
	switch field.Node.ShortTag() {
	case "!!null", "!vault", "!env", "!secret":
		return nil
	}
	value := strings.TrimSpace(field.Node.Value)
	if value == "" || isSecretReference(value) {
		return nil
	}
	return []entity.Problem{{
		RuleID:         "plaintext-password",
		Path:           field.Path,
		Severity:       entity.High,
		Explanation:    "Пароль записан в конфигурации открытым значением.",
		Recommendation: "Используйте ссылку на секрет или зашифрованное хранилище, например Ansible Vault.",
	}}
}

func isSecretReference(value string) bool {
	var expression string
	switch {
	case strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}"):
		expression = value[2 : len(value)-1]
	case strings.HasPrefix(value, "{{") && strings.HasSuffix(value, "}}"):
		expression = value[2 : len(value)-2]
	default:
		return false
	}
	return strings.TrimSpace(expression) != "" && !strings.ContainsAny(expression, "{}\r\n")
}
