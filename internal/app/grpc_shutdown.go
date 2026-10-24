package app

import (
	"log"
	"time"

	"google.golang.org/grpc"
)

func shutdownGRPC(server *grpc.Server, timeout time.Duration) {
	finished := make(chan struct{})
	go func() {
		server.GracefulStop()
		close(finished)
	}()
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case <-finished:
	case <-timer.C:
		log.Print("gRPC shutdown timeout reached; forcing stop")
		server.Stop()
	}
}
