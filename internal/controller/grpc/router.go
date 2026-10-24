package grpc

import (
	grpcv1 "github.com/MaksimovYuriy/ConfigAnalyzer/docs/proto/v1"
	v1 "github.com/MaksimovYuriy/ConfigAnalyzer/internal/controller/grpc/v1"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/usecase"
	pbgrpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewRouter(server *pbgrpc.Server, configuration usecase.Configuration) {
	controller := v1.NewController(configuration)
	grpcv1.RegisterConfigServiceServer(server, controller)
	reflection.Register(server)
}
