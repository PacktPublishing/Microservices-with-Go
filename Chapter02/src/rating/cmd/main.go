package main

import (
	"log"
	"net/http"

	"movieexample.com/rating/internal/controller"
	httphandler "movieexample.com/rating/internal/handler/http"
	"movieexample.com/rating/internal/repository/memory"
)

func main() {
	log.Println("Starting the rating service")
	repo := memory.New()
	svc := controller.New(repo)
	h := httphandler.New(svc)
	http.Handle("/rating", http.HandlerFunc(h.Handle))
	http.ListenAndServe(":8082", nil)
}
