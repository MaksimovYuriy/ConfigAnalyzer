package restapi

import "github.com/MaksimovYuriy/ConfigAnalyzer/internal/usecase"

type Controller struct {
	configuration usecase.Configuration
}

func NewController(configuration usecase.Configuration) *Controller {
	return &Controller{configuration: configuration}
}
