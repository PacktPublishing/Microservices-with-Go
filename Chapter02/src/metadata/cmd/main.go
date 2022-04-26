package main

import (
	"log"
	"context"
	"net/http"

	"movieexample.com/metadata/internal/controller"
	httphandler "movieexample.com/metadata/internal/handler/http"
	"movieexample.com/metadata/internal/repository/memory"
	"movieexample.com/metadata/pkg/model"
)

func main() {
	log.Println("Starting the movie metadata service")
	repo := memory.New()
	repo.Put(context.Background(), "tst", &model.Metadata{ID: "tst"})
	svc := controller.New(repo)
	h := httphandler.New(svc)
	http.Handle("/metadata", http.HandlerFunc(h.GetMetadataByID))
	http.ListenAndServe(":8081", nil)
}
