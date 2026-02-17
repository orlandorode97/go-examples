package main

import (
	"context"
	"log"
	"net"

	"github.com/orlandorode97/grpc-mesh-example/generated/service_c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	service_c.UnimplementedServiceCServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Ping(ctx context.Context, req *service_c.PingRequest) (*service_c.PingResponse, error) {
	return &service_c.PingResponse{Message: "C -> pong"}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	server := NewServer()
	service_c.RegisterServiceCServer(grpcServer, server)
	reflection.Register(grpcServer)
	log.Println("Service C running on :50051")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
