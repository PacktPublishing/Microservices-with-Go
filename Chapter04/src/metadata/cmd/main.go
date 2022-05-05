package main

import (
	"log"
	"net/http"

	"movieexample.com/metadata/internal/controller"
	httphandler "movieexample.com/metadata/internal/handler/http"
	"movieexample.com/metadata/internal/repository/memory"
)

func main() {
	log.Println("Starting the movie metadata service")
	repo := memory.New()
	svc := controller.New(repo)
	h := httphandler.New(svc)
	http.Handle("/metadata", http.HandlerFunc(h.GetMetadataByID))
	http.ListenAndServe(":8081", nil)
}
