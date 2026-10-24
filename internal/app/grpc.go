package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/config"
	router "github.com/MaksimovYuriy/ConfigAnalyzer/internal/controller/grpc"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/usecase"
	"google.golang.org/grpc"
)

func serveGRPC(ctx context.Context, cfg config.GRPCConfig, configurationUseCase usecase.Configuration) error {
	listener, err := net.Listen("tcp", net.JoinHostPort(cfg.Address, cfg.Port))
	if err != nil {
		return fmt.Errorf("gRPC listen: %w", err)
	}
	server := grpc.NewServer(grpc.MaxRecvMsgSize(cfg.MaxRecvMsgSize))
	router.NewRouter(server, configurationUseCase)
	log.Printf("gRPC listening on %s", listener.Addr())
	return serveGRPCServer(ctx, listener, server, cfg.ShutdownTimeout)
}

func serveGRPCServer(ctx context.Context, listener net.Listener, server *grpc.Server, timeout time.Duration) error {
	defer server.Stop()
	finished := make(chan error, 1)
	go func() {
		finished <- server.Serve(listener)
	}()
	select {
	case err := <-finished:
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			return fmt.Errorf("gRPC serve: %w", err)
		}
		return nil
	case <-ctx.Done():
		shutdownGRPC(server, timeout)
		return nil
	}
}
