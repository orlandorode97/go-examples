package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/orlandorode97/grpc-mesh-example/generated/service_b"
	"github.com/orlandorode97/grpc-mesh-example/generated/service_c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	service_b.UnimplementedServiceBServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Ping(ctx context.Context, req *service_b.PingRequest) (*service_b.PingResponse, error) {
	conn, err := grpc.NewClient(os.Getenv("SERVICE_C_ADDR"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	serverC := service_c.NewServiceCClient(conn)
	resp, err := serverC.Ping(ctx, &service_c.PingRequest{})
	if err != nil {
		return nil, err
	}
	return &service_b.PingResponse{Message: "B -> " + resp.Message}, nil

}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	server := NewServer()
	service_b.RegisterServiceBServer(grpcServer, server)
	reflection.Register(grpcServer)
	log.Println("Service B running on :50051")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
