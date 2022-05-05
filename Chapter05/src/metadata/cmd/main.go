package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"movieexample.com/gen"
	"movieexample.com/metadata/internal/controller"
	grpchandler "movieexample.com/metadata/internal/handler/grpc"
	"movieexample.com/metadata/internal/repository/memory"
)

func main() {
	log.Println("Starting the movie metadata service")
	repo := memory.New()
	svc := controller.New(repo)
	h := grpchandler.New(svc)
	lis, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	gen.RegisterMetadataServiceServer(srv, h)
	srv.Serve(lis)
}
