package main

import (
	"context"
	"log"
	"net/http"

	"movieexample.com/movie/internal/controller"
	metadatagateway "movieexample.com/movie/internal/gateway/metadata/http"
	ratinggateway "movieexample.com/movie/internal/gateway/rating/http"
	httphandler "movieexample.com/movie/internal/handler/http"
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
	h := httphandler.New(svc)
	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	http.ListenAndServe(":8083", nil)
}
