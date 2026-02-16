package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc/reflection"

	hellopb "github.com/orlandorode97/grpc-server/generated"
	proto "github.com/orlandorode97/grpc-server/generated"
	"google.golang.org/grpc"
)

type Server struct {
	hellopb.UnimplementedHelloServiceServer
}

func (s *Server) SayHello(ctx context.Context, req *hellopb.HelloRequest) (*hellopb.HelloResponse, error) {
	return &hellopb.HelloResponse{Message: "Hello, " + req.Name}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	helloService := Server{}
	proto.RegisterHelloServiceServer(grpcServer, &helloService)

	reflection.Register(grpcServer)

	log.Printf("server started on 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
