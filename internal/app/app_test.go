package app

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/config"
)

func testConfig() config.Config {
	return config.Config{
		HTTP: config.HTTPConfig{Address: "127.0.0.1", Port: "0", ShutdownTimeout: time.Second},
		GRPC: config.GRPCConfig{Address: "127.0.0.1", Port: "0", MaxRecvMsgSize: 1 << 20, ShutdownTimeout: time.Second},
	}
}

func TestAppStopsBothServersOnCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	finished := make(chan error, 1)
	go func() { finished <- run(ctx, testConfig()) }()
	waitStopped(t, finished)
}

func TestAppStopsWhenEitherPortIsOccupied(t *testing.T) {
	for _, transport := range []string{"HTTP", "gRPC"} {
		t.Run(transport, func(t *testing.T) {
			listener, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}
			defer listener.Close()
			_, port, err := net.SplitHostPort(listener.Addr().String())
			if err != nil {
				t.Fatal(err)
			}
			cfg := testConfig()
			if transport == "HTTP" {
				cfg.HTTP.Port = port
			} else {
				cfg.GRPC.Port = port
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := run(ctx, cfg); err == nil {
				t.Fatal("expected occupied port error")
			}
			if ctx.Err() != nil {
				t.Fatal("application did not stop the other server promptly")
			}
		})
	}
}
