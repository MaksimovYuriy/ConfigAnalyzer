package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/config"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/controller/restapi"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/usecase"
)

func serveHTTP(ctx context.Context, cfg config.HTTPConfig, configurationUseCase usecase.Configuration) error {
	listener, err := net.Listen("tcp", net.JoinHostPort(cfg.Address, cfg.Port))
	if err != nil {
		return fmt.Errorf("HTTP listen: %w", err)
	}
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
	log.Printf("HTTP listening on http://%s", listener.Addr())

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
