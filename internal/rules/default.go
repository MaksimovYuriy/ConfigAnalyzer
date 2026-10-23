package rules

import (
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/analyzer"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/rules/algorithm"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/rules/listen"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/rules/logging"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/rules/password"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/rules/tls"
)

func Default() []analyzer.Rule {
	return []analyzer.Rule{logging.Check, password.Check, tls.Check, listen.Check, algorithm.Check}
}
