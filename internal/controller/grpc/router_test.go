package grpc_test

import (
	"testing"

	router "github.com/MaksimovYuriy/ConfigAnalyzer/internal/controller/grpc"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/rules"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/usecase/configuration"
	"google.golang.org/grpc"
)

func TestNewRouterRegistersServices(t *testing.T) {
	server := grpc.NewServer()
	defer server.Stop()
	router.NewRouter(server, configuration.New(rules.Default()...))

	services := server.GetServiceInfo()
	service, ok := services["grpc.v1.ConfigService"]
	if !ok {
		t.Fatal("ConfigService was not registered")
	}
	if len(service.Methods) != 1 || service.Methods[0].Name != "Analyze" || service.Methods[0].IsClientStream || service.Methods[0].IsServerStream {
		t.Fatalf("unexpected service methods: %+v", service.Methods)
	}
	for _, name := range []string{"grpc.reflection.v1.ServerReflection", "grpc.reflection.v1alpha.ServerReflection"} {
		if _, ok := services[name]; !ok {
			t.Errorf("reflection service %s was not registered", name)
		}
	}
}
