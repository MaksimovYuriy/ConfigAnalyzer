package configuration

import (
	"context"
	"fmt"
	"strings"

	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/analyzer"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/entity"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/usecase"
	"github.com/MaksimovYuriy/ConfigAnalyzer/pkg/parser"
)

var _ usecase.Configuration = (*UseCase)(nil)

type UseCase struct {
	rules entity.RuleSet
}

func New(rules entity.RuleSet) *UseCase {
	return &UseCase{rules: entity.RuleSet{
		Fields: append([]entity.FieldRule(nil), rules.Fields...),
		Files:  append([]entity.FileRule(nil), rules.Files...),
	}}
}

func (uc *UseCase) Analyze(ctx context.Context, configuration entity.Configuration) ([]entity.Problem, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	document, err := parser.Parse(strings.NewReader(configuration.Content))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", usecase.ErrInvalidConfiguration, err)
	}
	problems := analyzer.Analyze(document.Root, configuration.FileMetadata, uc.rules)
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return problems, nil
}
