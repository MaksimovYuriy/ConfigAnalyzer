package rules

import (
	"strings"
	"unicode"
)

// nameSet matches explicitly registered aliases after normalizing spelling.
// It deliberately does not perform substring matching or infer synonyms.
type nameSet map[string]struct{}

func newNameSet(aliases ...string) nameSet {
	names := make(nameSet, len(aliases))
	for _, alias := range aliases {
		if name := normalizeName(alias); name != "" {
			names[name] = struct{}{}
		}
	}
	return names
}

func (names nameSet) matches(name string) bool {
	_, ok := names[normalizeName(name)]
	return ok
}

// normalizeName preserves punctuation other than separators used in common
// parameter spellings. In particular, dots remain meaningful path separators.
func normalizeName(name string) string {
	return strings.Map(func(r rune) rune {
		if r == '-' || r == '_' {
			return -1
		}
		return unicode.ToLower(r)
	}, name)
}

// These sets identify parameter names, not vulnerabilities. Rules must also
// consider the value and, where necessary, the surrounding configuration.
var (
	passwordNames = newNameSet("password", "passwd", "pwd", "db_password", "database_password")

	loggingNames  = newNameSet("log", "logs", "logging", "logger")
	logLevelNames = newNameSet("level", "log_level", "logging_level")

	// Keep opposite boolean meanings separate when implementing TLS checks.
	tlsEnabledNames    = newNameSet("tls_enabled", "ssl_enabled")
	tlsVerifyNames     = newNameSet("tls_verify", "ssl_verify", "verify_ssl", "verify_tls", "verify_certificate", "verify_certificates")
	tlsSkipVerifyNames = newNameSet("insecure_skip_verify", "skip_tls_verify", "skip_ssl_verify")

	listenAddressNames = newNameSet("host", "bind", "bind_address", "listen", "listen_address", "address")
	algorithmNames     = newNameSet("algorithm", "digest_algorithm", "hash_algorithm", "hashing_algorithm", "cipher")
)
