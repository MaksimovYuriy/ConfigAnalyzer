package app

import (
	"context"
	"net"
	"testing"

	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/config"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/usecase/configuration"
)

func TestHTTPStopsOnCancellation(t *testing.T) {
	t.Setenv("HTTP_ADDRESS", "127.0.0.1")
	t.Setenv("HTTP_PORT", "8081")
	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	cfg.HTTP.Port = "0"
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := serveHTTP(ctx, cfg.HTTP, configuration.New()); err != nil {
		t.Fatal(err)
	}
}

func TestHTTPReturnsListenError(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	host, port, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	err = serveHTTP(context.Background(), config.HTTPConfig{Address: host, Port: port}, configuration.New())
	if err == nil {
		t.Fatal("expected occupied address error")
	}
}
