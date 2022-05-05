package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"movieexample.com/gen"
	"movieexample.com/metadata/internal/controller"
	"movieexample.com/metadata/internal/repository"
	"movieexample.com/metadata/pkg/model"
)

// Handler defines a movie metadata gRPC handler.
type Handler struct {
	gen.UnimplementedMetadataServiceServer
	svc *controller.MetadataService
}

// New creates a new movie metadata gRPC handler.
func New(svc *controller.MetadataService) *Handler {
	return &Handler{svc: svc}
}

// GetMetadataByID returns movie metadata by id.
func (h *Handler) GetMetadataByID(ctx context.Context, req *gen.GetMetadataByIDRequest) (*gen.GetMetadataByIDResponse, error) {
	if req == nil || req.MovieId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}
	m, err := h.svc.Get(ctx, req.MovieId)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &gen.GetMetadataByIDResponse{Metadata: model.MetadataToProto(m)}, nil
}
