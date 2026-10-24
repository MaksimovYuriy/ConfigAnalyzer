package configuration_test

import (
	"context"
	"errors"
	"testing"

	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/entity"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/rules"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/usecase"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/usecase/configuration"
)

func TestConfigurationAnalyze(t *testing.T) {
	service := configuration.New(rules.Default()...)
	for _, tt := range []struct {
		name, content string
		findings      int
		invalid       bool
	}{
		{"YAML", "log: {level: debug}", 1, false},
		{"JSON", `{"password":"secret","tls_enabled":false}`, 2, false},
		{"safe", "log: {level: info}", 0, false},
		{"scalar", "5", 0, false},
		{"invalid", "log: [", 0, true},
		{"empty", "", 0, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			problems, err := service.Analyze(context.Background(), entity.Configuration{Content: tt.content})
			if tt.invalid {
				if !errors.Is(err, usecase.ErrInvalidConfiguration) {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err != nil || len(problems) != tt.findings {
				t.Fatalf("problems=%+v, error=%v", problems, err)
			}
		})
	}
}

func TestConfigurationCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := configuration.New().Analyze(ctx, entity.Configuration{Content: "5"})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error=%v", err)
	}
}

func TestConfigurationCustomRule(t *testing.T) {
	custom := func(field entity.Field) []entity.Problem {
		if field.Key == "custom" {
			return []entity.Problem{{RuleID: "custom", Path: field.Path}}
		}
		return nil
	}
	service := configuration.New(custom)
	problems, err := service.Analyze(context.Background(), entity.Configuration{Content: "custom: true"})
	if err != nil || len(problems) != 1 || problems[0].RuleID != "custom" {
		t.Fatalf("problems=%+v, error=%v", problems, err)
	}
}
