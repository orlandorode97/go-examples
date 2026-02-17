package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/orlandorode97/grpc-mesh-example/generated/service_a"
	"github.com/orlandorode97/grpc-mesh-example/generated/service_b"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	service_a.UnimplementedServiceAServer
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Ping(ctx context.Context, req *service_a.PingRequest) (*service_a.PingResponse, error) {
	conn, err := grpc.NewClient(os.Getenv("SERVICE_B_ADDR"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	serverB := service_b.NewServiceBClient(conn)
	resp, err := serverB.Ping(ctx, &service_b.PingRequest{})
	if err != nil {
		return nil, err
	}
	return &service_a.PingResponse{Message: "A -> " + resp.Message}, nil

}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	server := NewServer()
	service_a.RegisterServiceAServer(grpcServer, server)
	reflection.Register(grpcServer)
	log.Println("Service A running on :50051")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
