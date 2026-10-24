package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTP HTTPConfig `env-prefix:"HTTP_"`
}

type HTTPConfig struct {
	Port              string        `env:"PORT" env-default:"8081"`
	Address           string        `env:"ADDRESS" env-default:"0.0.0.0"`
	ReadTimeout       time.Duration `env:"READ_TIMEOUT" env-default:"60s"`
	WriteTimeout      time.Duration `env:"WRITE_TIMEOUT" env-default:"60s"`
	ReadHeaderTimeout time.Duration `env:"READ_HEADER_TIMEOUT" env-default:"5s"`
	IdleTimeout       time.Duration `env:"IDLE_TIMEOUT" env-default:"60s"`
	ShutdownTimeout   time.Duration `env:"SHUTDOWN_TIMEOUT" env-default:"10s"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("read .env: %w", err)
		}
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			return nil, fmt.Errorf("read environment: %w", err)
		}
	}
	if err := cfg.HTTP.validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (cfg HTTPConfig) validate() error {
	port, err := strconv.Atoi(cfg.Port)
	if err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("HTTP_PORT must be a number from 1 to 65535")
	}
	for name, duration := range map[string]time.Duration{
		"HTTP_READ_TIMEOUT":        cfg.ReadTimeout,
		"HTTP_WRITE_TIMEOUT":       cfg.WriteTimeout,
		"HTTP_READ_HEADER_TIMEOUT": cfg.ReadHeaderTimeout,
		"HTTP_IDLE_TIMEOUT":        cfg.IdleTimeout,
		"HTTP_SHUTDOWN_TIMEOUT":    cfg.ShutdownTimeout,
	} {
		if duration <= 0 {
			return fmt.Errorf("%s must be positive", name)
		}
	}
	return nil
}
