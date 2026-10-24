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
	service := configuration.New(rules.Default())
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
	_, err := configuration.New(entity.RuleSet{}).Analyze(ctx, entity.Configuration{Content: "5"})
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
	service := configuration.New(entity.RuleSet{Fields: []entity.FieldRule{custom}})
	problems, err := service.Analyze(context.Background(), entity.Configuration{Content: "custom: true"})
	if err != nil || len(problems) != 1 || problems[0].RuleID != "custom" {
		t.Fatalf("problems=%+v, error=%v", problems, err)
	}
}

func TestConfigurationFileMetadata(t *testing.T) {
	service := configuration.New(rules.Default())
	input := entity.Configuration{
		Content:      "5",
		FileMetadata: &entity.FileMetadata{Path: "config.yaml", Mode: 0666},
	}
	problems, err := service.Analyze(context.Background(), input)
	if err != nil || len(problems) != 1 || problems[0].RuleID != "world-writable-config" {
		t.Fatalf("problems=%+v, error=%v", problems, err)
	}
	input.FileMetadata = nil
	problems, err = service.Analyze(context.Background(), input)
	if err != nil || len(problems) != 0 {
		t.Fatalf("without metadata: problems=%+v, error=%v", problems, err)
	}
}

func TestConfigurationCopiesRuleSet(t *testing.T) {
	checks := rules.Default()
	service := configuration.New(checks)
	checks.Files[0] = func(entity.FileMetadata) []entity.Problem { return nil }
	checks.Fields[0] = func(entity.Field) []entity.Problem { return nil }
	problems, err := service.Analyze(context.Background(), entity.Configuration{
		Content:      "log: {level: debug}",
		FileMetadata: &entity.FileMetadata{Path: "config.yaml", Mode: 0666},
	})
	if err != nil || len(problems) != 2 {
		t.Fatalf("rules were changed: problems=%+v, error=%v", problems, err)
	}
}
