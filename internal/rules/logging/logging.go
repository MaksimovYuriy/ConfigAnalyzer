package logging

import (
	"strings"

	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/entity"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/rules/helpers"
)

var (
	loggingNames = helpers.NewNames(
		"log",
		"logs",
		"logging",
		"logger",
	)
	logLevelNames = helpers.NewNames("level", "log_level", "logging_level")
)

func Check(field entity.Field) []entity.Problem {
	if !logLevelNames.Matches(field.Key) {
		return nil
	}
	value, isString := helpers.ScalarText(field)
	if !isString || !strings.EqualFold(value, "debug") {
		return nil
	}
	if helpers.Normalize(field.Key) == "level" && !helpers.InContext(field, loggingNames) {
		return nil
	}
	return []entity.Problem{{
		RuleID:         "debug-logging",
		Path:           field.Path,
		Severity:       entity.Low,
		Explanation:    "Логирование в debug-режиме.",
		Recommendation: "Поменяйте уровень логирования на info или выше.",
	}}
}
