package config

import (
	"os"
	"strings"
	"testing"
	"time"
)

var httpEnvironment = []string{
	"HTTP_ADDRESS", "HTTP_PORT", "HTTP_READ_TIMEOUT", "HTTP_WRITE_TIMEOUT",
	"HTTP_READ_HEADER_TIMEOUT", "HTTP_IDLE_TIMEOUT", "HTTP_SHUTDOWN_TIMEOUT",
}

func clearEnvironment(t *testing.T) {
	t.Helper()
	t.Chdir(t.TempDir())
	for _, name := range httpEnvironment {
		t.Setenv(name, "")
		if err := os.Unsetenv(name); err != nil {
			t.Fatal(err)
		}
	}
}

func TestLoadDefaults(t *testing.T) {
	clearEnvironment(t)
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	want := HTTPConfig{
		Address: "0.0.0.0", Port: "8081", ReadTimeout: 60 * time.Second,
		WriteTimeout: 60 * time.Second, ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout: 60 * time.Second, ShutdownTimeout: 10 * time.Second,
	}
	if cfg.HTTP != want {
		t.Fatalf("HTTP config=%+v, want=%+v", cfg.HTTP, want)
	}
}

func TestLoadEnvironment(t *testing.T) {
	clearEnvironment(t)
	t.Setenv("HTTP_ADDRESS", "::1")
	t.Setenv("HTTP_PORT", "9090")
	for _, name := range httpEnvironment[2:] {
		t.Setenv(name, "250ms")
	}
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	want := HTTPConfig{
		Address: "::1", Port: "9090", ReadTimeout: 250 * time.Millisecond,
		WriteTimeout: 250 * time.Millisecond, ReadHeaderTimeout: 250 * time.Millisecond,
		IdleTimeout: 250 * time.Millisecond, ShutdownTimeout: 250 * time.Millisecond,
	}
	if cfg.HTTP != want {
		t.Fatalf("HTTP config=%+v, want=%+v", cfg.HTTP, want)
	}
}

func TestLoadInvalidPort(t *testing.T) {
	for _, port := range []string{"0", "-1", "65536", "http"} {
		t.Run(port, func(t *testing.T) {
			clearEnvironment(t)
			t.Setenv("HTTP_PORT", port)
			if _, err := Load(); err == nil {
				t.Fatal("expected port error")
			}
		})
	}
}

func TestLoadInvalidTimeout(t *testing.T) {
	for _, name := range httpEnvironment[2:] {
		for _, value := range []string{"invalid", "0s", "-1s", "999999999999999999999h"} {
			t.Run(name+value, func(t *testing.T) {
				clearEnvironment(t)
				t.Setenv(name, value)
				_, err := Load()
				if err == nil || !strings.Contains(err.Error(), name) {
					t.Fatalf("expected error naming %s, got %v", name, err)
				}
			})
		}
	}
}

func TestLoadEnvFile(t *testing.T) {
	clearEnvironment(t)
	t.Setenv("HTTP_PORT", "8088")
	data := []byte("HTTP_ADDRESS=127.0.0.1\nHTTP_PORT=9091\nHTTP_READ_TIMEOUT=2s\n")
	if err := os.WriteFile(".env", data, 0600); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.HTTP.Address != "127.0.0.1" || cfg.HTTP.Port != "9091" || cfg.HTTP.ReadTimeout != 2*time.Second {
		t.Fatalf("unexpected .env config: %+v", cfg.HTTP)
	}
}

func TestLoadInvalidEnvFile(t *testing.T) {
	clearEnvironment(t)
	if err := os.WriteFile(".env", []byte("HTTP_PORT='unterminated"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(); err == nil {
		t.Fatal("expected .env parse error")
	}
}
