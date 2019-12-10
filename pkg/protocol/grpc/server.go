package grpc

import (
	"context"
	"net"
	"os"
	"os/signal"

	"github.com/arrowfeng/go-grpc-http-rest-microservice-demo/pkg/logger"
	"github.com/arrowfeng/go-grpc-http-rest-microservice-demo/pkg/protocol/grpc/middleware"

	"google.golang.org/grpc"

	v1 "github.com/arrowfeng/go-grpc-http-rest-microservice-demo/pkg/api/v1"
)

// RunServer runs gRPC service to publish ToDo service
func RunServer(ctx context.Context, v1API v1.ToDoServiceServer, port string) error {

	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	// gRPC server statup options
	opts := []grpc.ServerOption{}

	// add middleware
	opts = middleware.AddLogging(logger.Log, opts)

	// register service
	server := grpc.NewServer(opts...)
	v1.RegisterToDoServiceServer(server, v1API)

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			// sig is a ^C, handle it
			logger.Log.Warn("shutting down RPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	// start gRPC server
	logger.Log.Warn("starting gRPC server...")
	return server.Serve(listen)
}
