package rules

import (
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/entity"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/rules/algorithm"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/rules/listen"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/rules/logging"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/rules/password"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/rules/tls"
)

func Default() []entity.Rule {
	return []entity.Rule{logging.Check, password.Check, tls.Check, listen.Check, algorithm.Check}
}
