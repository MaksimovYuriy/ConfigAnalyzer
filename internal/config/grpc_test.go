package config

import (
	"os"
	"strings"
	"testing"
	"time"
)

var grpcEnvironment = []string{
	"GRPC_ADDRESS", "GRPC_PORT", "GRPC_MAX_RECV_MSG_SIZE", "GRPC_SHUTDOWN_TIMEOUT",
}

func TestLoadGRPCDefaults(t *testing.T) {
	clearEnvironment(t)
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	want := GRPCConfig{
		Address: "0.0.0.0", Port: "50051",
		MaxRecvMsgSize: 1 << 20, ShutdownTimeout: 10 * time.Second,
	}
	if cfg.GRPC != want {
		t.Fatalf("gRPC config=%+v, want=%+v", cfg.GRPC, want)
	}
}

func TestLoadGRPCEnvironment(t *testing.T) {
	clearEnvironment(t)
	t.Setenv("GRPC_ADDRESS", "::1")
	t.Setenv("GRPC_PORT", "50052")
	t.Setenv("GRPC_MAX_RECV_MSG_SIZE", "2097152")
	t.Setenv("GRPC_SHUTDOWN_TIMEOUT", "2s")
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	want := GRPCConfig{
		Address: "::1", Port: "50052",
		MaxRecvMsgSize: 2 << 20, ShutdownTimeout: 2 * time.Second,
	}
	if cfg.GRPC != want {
		t.Fatalf("gRPC config=%+v, want=%+v", cfg.GRPC, want)
	}
	if cfg.HTTP.Port != "8081" {
		t.Fatal("gRPC settings changed HTTP configuration")
	}
}

func TestLoadInvalidGRPCSettings(t *testing.T) {
	for _, tt := range []struct {
		name, value string
	}{
		{"GRPC_PORT", "0"},
		{"GRPC_PORT", "-1"},
		{"GRPC_PORT", "65536"},
		{"GRPC_PORT", "invalid"},
		{"GRPC_MAX_RECV_MSG_SIZE", "0"},
		{"GRPC_MAX_RECV_MSG_SIZE", "-1"},
		{"GRPC_MAX_RECV_MSG_SIZE", "invalid"},
		{"GRPC_SHUTDOWN_TIMEOUT", "0s"},
		{"GRPC_SHUTDOWN_TIMEOUT", "-1s"},
		{"GRPC_SHUTDOWN_TIMEOUT", "invalid"},
	} {
		t.Run(tt.name+tt.value, func(t *testing.T) {
			clearEnvironment(t)
			t.Setenv(tt.name, tt.value)
			_, err := Load()
			if err == nil || !strings.Contains(err.Error(), tt.name) {
				t.Fatalf("expected error naming %s, got %v", tt.name, err)
			}
		})
	}
}

func TestLoadGRPCEnvFile(t *testing.T) {
	clearEnvironment(t)
	data := []byte("GRPC_ADDRESS=127.0.0.1\nGRPC_PORT=50053\nGRPC_MAX_RECV_MSG_SIZE=2048\nGRPC_SHUTDOWN_TIMEOUT=500ms\n")
	if err := os.WriteFile(".env", data, 0600); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	want := GRPCConfig{
		Address: "127.0.0.1", Port: "50053",
		MaxRecvMsgSize: 2048, ShutdownTimeout: 500 * time.Millisecond,
	}
	if cfg.GRPC != want {
		t.Fatalf("gRPC config=%+v, want=%+v", cfg.GRPC, want)
	}
}
