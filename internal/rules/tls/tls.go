package tls

import (
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/analyzer"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/rules/helpers"
)

var (
	tlsNames     = helpers.NewNames("tls", "ssl")
	enabledNames = helpers.NewNames("enabled", "enable")
	verifyNames  = helpers.NewNames(
		"verify",
		"verify_peer",
		"verify_certificate",
		"verify_certificates",
	)
	disableTLSNames = helpers.NewNames(
		"disable_tls",
		"disable_ssl",
		"tls_disabled",
		"ssl_disabled",
	)
	tlsEnabledNames = helpers.NewNames("tls_enabled", "ssl_enabled")
	tlsVerifyNames  = helpers.NewNames(
		"tls_verify",
		"ssl_verify",
		"verify_ssl",
		"verify_tls",
		"verify_certificate",
		"verify_certificates",
		"verify_peer",
		"verify_certs",
		"reject_unauthorized",
		"tls_verify_certificate",
		"ssl_verify_certificate",
	)
	tlsSkipVerifyNames = helpers.NewNames(
		"insecure_skip_verify",
		"skip_tls_verify",
		"skip_ssl_verify",
		"skip_verify",
	)
)

func Check(field analyzer.Field) []analyzer.Problem {
	value, isBoolean := helpers.BooleanValue(field)
	if !isBoolean {
		return nil
	}
	inTLS := helpers.InContext(field, tlsNames)
	enablesTLS := tlsNames.Matches(field.Key) || tlsEnabledNames.Matches(field.Key) || (inTLS && enabledNames.Matches(field.Key))
	verifiesTLS := tlsVerifyNames.Matches(field.Key) || (inTLS && verifyNames.Matches(field.Key))
	disabled := (!value && enablesTLS) || (value && disableTLSNames.Matches(field.Key))
	unverified := (value && tlsSkipVerifyNames.Matches(field.Key)) || (!value && verifiesTLS)
	if disabled {
		return []analyzer.Problem{{
			RuleID:         "tls-disabled",
			Path:           field.Path,
			Severity:       analyzer.High,
			Explanation:    "TLS отключён.",
			Recommendation: "Включите TLS для защиты соединения.",
		}}
	}
	if unverified {
		return []analyzer.Problem{{
			RuleID:         "tls-verification-disabled",
			Path:           field.Path,
			Severity:       analyzer.High,
			Explanation:    "Проверка TLS-сертификатов отключена.",
			Recommendation: "Включите проверку сертификата и доверенной цепочки центра сертификации.",
		}}
	}
	return nil
}
