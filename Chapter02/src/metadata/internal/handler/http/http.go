package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"movieexample.com/metadata/internal/controller"
	"movieexample.com/metadata/internal/repository"
)

// Handler defines a movie metadata HTTP handler.
type Handler struct {
	svc *controller.MetadataService
}

// New creates a new movie metadata HTTP handler.
func New(svc *controller.MetadataService) *Handler {
	return &Handler{svc}
}

// GetMetadataByID handles GET /metadata requests.
func (h *Handler) GetMetadataByID(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx := req.Context()
	m, err := h.svc.Get(ctx, id)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(m); err != nil {
		log.Printf("Response encode error: %v\n", err)
	}
}
