package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/config"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/controller/restapi"
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
	return serve(ctx, cfg.HTTP)
}

func serve(ctx context.Context, cfg config.HTTPConfig) error {
	listener, err := net.Listen("tcp", net.JoinHostPort(cfg.Address, cfg.Port))
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	configurationUseCase := configuration.New(rules.Default()...)
	controller := restapi.NewController(configurationUseCase)
	server := &http.Server{
		Handler:           restapi.NewRouter(controller),
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}
	finished := make(chan error, 1)
	go func() {
		finished <- server.Serve(listener)
	}()
	log.Printf("API listening on http://%s", listener.Addr())

	select {
	case err := <-finished:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			server.Close()
			return fmt.Errorf("shutdown: %w", err)
		}
		return nil
	}
}
