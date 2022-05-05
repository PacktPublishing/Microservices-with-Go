package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"movieexample.com/gen"
	"movieexample.com/movie/internal/controller"
	metadatagateway "movieexample.com/movie/internal/gateway/metadata/grpc"
	ratinggateway "movieexample.com/movie/internal/gateway/rating/grpc"
	grpchandler "movieexample.com/movie/internal/handler/grpc"
	"movieexample.com/pkg/discovery/static"
)

func main() {
	log.Println("Starting the movie service")
	registry := static.NewRegistry(map[string][]string{
		"metadata": {"localhost:8081"},
		"rating":   {"localhost:8082"},
		"movie":    {"localhost:8083"},
	})
	ctx := context.Background()
	if err := registry.Register(ctx, "movie", "localhost:8083"); err != nil {
		panic(err)
	}
	defer registry.Deregister(ctx, "movie")
	metadataGateway := metadatagateway.New(registry)
	ratingGateway := ratinggateway.New(registry)
	svc := controller.New(ratingGateway, metadataGateway)
	h := grpchandler.New(svc)
	lis, err := net.Listen("tcp", "localhost:8083")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	gen.RegisterMovieServiceServer(srv, h)
	srv.Serve(lis)
}
