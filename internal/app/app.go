package app

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/config"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/rules"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/usecase/configuration"
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return run(ctx, *cfg)
}

func run(ctx context.Context, cfg config.Config) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	configurationUseCase := configuration.New(rules.Default()...)
	finished := make(chan error, 2)
	go func() {
		finished <- serveHTTP(ctx, cfg.HTTP, configurationUseCase)
	}()
	go func() {
		finished <- serveGRPC(ctx, cfg.GRPC, configurationUseCase)
	}()

	first := <-finished
	cancel()
	second := <-finished
	return errors.Join(first, second)
}
