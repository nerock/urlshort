package grpc

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
)

// Service is a service that can be registered in the argument provided grpc.Server
type Service interface {
	Register(*grpc.Server)
}

// Server is a gRPC server
type Server struct {
	srv *grpc.Server
}

// NewGRPCServer creates a new Server and registers all the provided Service
func NewGRPCServer(services ...Service) Server {
	srv := grpc.NewServer()
	for _, svc := range services {
		svc.Register(srv)
	}

	return Server{
		srv: srv,
	}
}

// RunServer starts the Server in the provided port
func (g Server) RunServer(port int, ) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %s", err)
	}

	return g.srv.Serve(lis)
}

// Shutdown gracefully shuts down Server
func (g Server) Shutdown() {
	g.srv.GracefulStop()
}
