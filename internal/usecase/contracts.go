package usecase

import (
	"context"

	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/entity"
)

type Configuration interface {
	Analyze(context.Context, entity.Configuration) ([]entity.Problem, error)
}
