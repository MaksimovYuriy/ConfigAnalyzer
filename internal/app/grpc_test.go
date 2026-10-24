package app

import (
	"context"
	"net"
	"testing"
	"time"

	grpcv1 "github.com/MaksimovYuriy/ConfigAnalyzer/docs/proto/v1"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/config"
	router "github.com/MaksimovYuriy/ConfigAnalyzer/internal/controller/grpc"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/entity"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/rules"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/usecase"
	"github.com/MaksimovYuriy/ConfigAnalyzer/internal/usecase/configuration"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func startServer(t *testing.T, service usecase.Configuration, timeout time.Duration, options ...grpc.ServerOption) (grpcv1.ConfigServiceClient, context.CancelFunc, <-chan error) {
	t.Helper()
	listener := bufconn.Listen(1 << 20)
	server := grpc.NewServer(options...)
	router.NewRouter(server, service)
	ctx, cancel := context.WithCancel(context.Background())
	finished := make(chan error, 1)
	go func() { finished <- serveGRPCServer(ctx, listener, server, timeout) }()
	connection, err := grpc.NewClient("passthrough:///test",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return listener.DialContext(ctx)
		}),
	)
	if err != nil {
		cancel()
		t.Fatal(err)
	}
	t.Cleanup(func() {
		connection.Close()
		cancel()
		server.Stop()
	})
	return grpcv1.NewConfigServiceClient(connection), cancel, finished
}

func waitStopped(t *testing.T, finished <-chan error) {
	t.Helper()
	select {
	case err := <-finished:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("server did not stop")
	}
}

func TestServerAnalyzeAndShutdown(t *testing.T) {
	client, cancel, finished := startServer(t, configuration.New(rules.Default()), time.Second)
	ctx, stop := context.WithTimeout(context.Background(), 5*time.Second)
	defer stop()
	response, err := client.Analyze(ctx, &grpcv1.AnalyzeRequest{Content: "logging: {level: debug}"})
	if err != nil {
		t.Fatal(err)
	}
	if len(response.Problems) != 1 || response.Problems[0].Severity != grpcv1.Severity_SEVERITY_LOW {
		t.Fatalf("unexpected response: %v", response)
	}
	_, err = client.Analyze(ctx, &grpcv1.AnalyzeRequest{Content: "log: ["})
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("error=%v", err)
	}
	cancel()
	waitStopped(t, finished)
}

type blockingConfiguration struct {
	started chan struct{}
	release chan struct{}
}

func (service *blockingConfiguration) Analyze(ctx context.Context, _ entity.Configuration) ([]entity.Problem, error) {
	close(service.started)
	select {
	case <-service.release:
		return nil, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func TestGracefulShutdownWaitsForRequest(t *testing.T) {
	service := &blockingConfiguration{started: make(chan struct{}), release: make(chan struct{})}
	client, cancel, finished := startServer(t, service, 5*time.Second)
	ctx, stop := context.WithTimeout(context.Background(), 5*time.Second)
	defer stop()
	requestDone := make(chan error, 1)
	go func() {
		_, err := client.Analyze(ctx, &grpcv1.AnalyzeRequest{Content: "5"})
		requestDone <- err
	}()
	select {
	case <-service.started:
	case <-ctx.Done():
		t.Fatal("request did not start")
	}
	cancel()
	select {
	case <-finished:
		t.Fatal("shutdown did not wait for the active request")
	case <-time.After(20 * time.Millisecond):
	}
	close(service.release)
	if err := <-requestDone; err != nil {
		t.Fatal(err)
	}
	waitStopped(t, finished)
}

func TestShutdownTimeoutCancelsRequest(t *testing.T) {
	service := &blockingConfiguration{started: make(chan struct{}), release: make(chan struct{})}
	client, cancel, finished := startServer(t, service, 20*time.Millisecond)
	ctx, stop := context.WithTimeout(context.Background(), 5*time.Second)
	defer stop()
	requestDone := make(chan error, 1)
	go func() {
		_, err := client.Analyze(ctx, &grpcv1.AnalyzeRequest{Content: "5"})
		requestDone <- err
	}()
	select {
	case <-service.started:
	case <-ctx.Done():
		t.Fatal("request did not start")
	}
	cancel()
	waitStopped(t, finished)
	if err := <-requestDone; err == nil {
		t.Fatal("forced stop did not interrupt the request")
	}
}

func TestServerMessageLimit(t *testing.T) {
	client, cancel, finished := startServer(t, configuration.New(entity.RuleSet{}), time.Second, grpc.MaxRecvMsgSize(16))
	ctx, stop := context.WithTimeout(context.Background(), 5*time.Second)
	defer stop()
	_, err := client.Analyze(ctx, &grpcv1.AnalyzeRequest{Content: "a configuration longer than sixteen bytes"})
	if status.Code(err) != codes.ResourceExhausted {
		t.Fatalf("error=%v", err)
	}
	cancel()
	waitStopped(t, finished)
}

func TestServeOccupiedPort(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	host, port, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	if err := serveGRPC(context.Background(), config.GRPCConfig{Address: host, Port: port}, configuration.New(entity.RuleSet{})); err == nil {
		t.Fatal("expected occupied port error")
	}
}
