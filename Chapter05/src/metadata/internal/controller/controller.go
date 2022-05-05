package controller

import (
	"context"

	"movieexample.com/metadata/pkg/model"
)

type metadataRepository interface {
	Get(ctx context.Context, id string) (*model.Metadata, error)
}

// MetadataService encapsulates the metadata service business logic.
type MetadataService struct {
	repo metadataRepository
}

// New creates a metadata service.
func New(repo metadataRepository) *MetadataService {
	return &MetadataService{repo}
}

// Get returns metadata by movie id.
func (s *MetadataService) Get(ctx context.Context, id string) (*model.Metadata, error) {
	return s.repo.Get(ctx, id)
}
