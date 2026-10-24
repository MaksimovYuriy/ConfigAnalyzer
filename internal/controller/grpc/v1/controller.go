package v1

import (
	grpcv1 "github.com/MaksimovYuriy/ConfigAnalyzer/docs/proto/v1"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/usecase"
)

var _ grpcv1.ConfigServiceServer = (*Controller)(nil)

type Controller struct {
	grpcv1.UnimplementedConfigServiceServer
	configuration usecase.Configuration
}

func NewController(configuration usecase.Configuration) *Controller {
	return &Controller{configuration: configuration}
}
