package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"movieexample.com/rating/internal/controller"
	"movieexample.com/rating/pkg/model"
)

// Handler defines a HTTP rating handler.
type Handler struct {
	svc *controller.RatingService
}

// New creates a new movie metadata HTTP handler.
func New(svc *controller.RatingService) *Handler {
	return &Handler{svc}
}

// Handle handles PUT and GET /rating requests.
func (h *Handler) Handle(w http.ResponseWriter, req *http.Request) {
	recordID := model.RecordID(req.FormValue("id"))
	if recordID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	typeInt, err := strconv.Atoi(req.FormValue("type"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	recordType := model.RecordType(typeInt)
	switch req.Method {
	case http.MethodGet:
		v, err := h.svc.GetAggregatedRating(req.Context(), recordID, recordType)
		if err != nil && errors.Is(err, controller.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err := json.NewEncoder(w).Encode(v); err != nil {
			log.Printf("Response encode error: %v\n", err)
		}
	case http.MethodPut:
		userID := model.UserID(req.FormValue("userId"))
		v, err := strconv.ParseFloat(req.FormValue("value"), 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := h.svc.PutRating(req.Context(), recordID, recordType, &model.Rating{UserID: userID, Value: model.RatingValue(v)}); err != nil {
			log.Printf("Repository put error: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
