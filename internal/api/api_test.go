package api

import (
	"context"
	"net"
	"testing"

	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/config"
)

func TestServeStopsOnCancellation(t *testing.T) {
	t.Setenv("HTTP_ADDRESS", "127.0.0.1")
	t.Setenv("HTTP_PORT", "8081")
	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	cfg.HTTP.Port = "0"
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := serve(ctx, cfg.HTTP); err != nil {
		t.Fatal(err)
	}
}

func TestServeReturnsListenError(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	host, port, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	err = serve(context.Background(), config.HTTPConfig{Address: host, Port: port})
	if err == nil {
		t.Fatal("expected occupied address error")
	}
}
