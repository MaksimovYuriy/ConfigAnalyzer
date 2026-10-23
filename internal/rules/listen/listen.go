package listen

import (
	"net"

	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/analyzer"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/rules/helpers"
)

var (
	listenAddressNames = helpers.NewNames(
		"host",
		"bind",
		"bind_address",
		"listen",
		"listen_address",
		"address",
	)
)

func Check(field analyzer.Field) []analyzer.Problem {
	if !helpers.NamedValue(field, listenAddressNames) {
		return nil
	}
	value, isString := helpers.ScalarText(field)
	if !isString {
		return nil
	}
	host := value
	if parsedHost, _, err := net.SplitHostPort(value); err == nil {
		host = parsedHost
	}
	ip := net.ParseIP(host)
	if ip == nil || ip.To4() == nil || !ip.IsUnspecified() {
		return nil
	}
	return []analyzer.Problem{{
		RuleID:         "wildcard-listen",
		Path:           field.Path,
		Severity:       analyzer.Medium,
		Explanation:    "Сервис слушает все IPv4-интерфейсы; ограничения доступа не установлены этой проверкой.",
		Recommendation: "Ограничьте адрес прослушивания или проверьте ограничения доступа на уровне сети.",
	}}
}
