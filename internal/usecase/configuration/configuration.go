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
	rules []entity.Rule
}

func New(rules ...entity.Rule) *UseCase {
	return &UseCase{rules: append([]entity.Rule(nil), rules...)}
}

func (uc *UseCase) Analyze(ctx context.Context, configuration entity.Configuration) ([]entity.Problem, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	document, err := parser.Parse(strings.NewReader(configuration.Content))
	if err != nil {
		return nil, fmt.Errorf("%w: %v", usecase.ErrInvalidConfiguration, err)
	}
	problems := analyzer.Analyze(document.Root, uc.rules...)
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return problems, nil
}
