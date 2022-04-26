package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"movieexample.com/movie/internal/controller"
)

// Handler defines a movie service.
type Handler struct {
	svc *controller.MovieService
}

// New creates a new movie HTTP handler.
func New(svc *controller.MovieService) *Handler {
	return &Handler{svc}
}

// GetMovieDetails handles GET /movie requests.
func (h *Handler) GetMovieDetails(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	details, err := h.svc.Get(req.Context(), id)
	if err != nil && errors.Is(err, controller.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(details); err != nil {
		log.Printf("Response encode error: %v\n", err)
	}
}
