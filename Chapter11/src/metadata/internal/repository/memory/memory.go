package memory

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel"
	"movieexample.com/metadata/internal/repository"
	"movieexample.com/metadata/pkg/model"
)

const tracerID = "metadata-repository-memory"

// Repository defines a memory movie matadata repository.
type Repository struct {
	sync.RWMutex
	data map[string]*model.Metadata
}

// New creates a new memory repository.
func New() *Repository {
	return &Repository{data: map[string]*model.Metadata{}}
}

// Get retrieves movie metadata for by movie id.
func (r *Repository) Get(ctx context.Context, id string) (*model.Metadata, error) {
	_, span := otel.Tracer(tracerID).Start(ctx, "Repository/Get")
	defer span.End()
	r.RLock()
	defer r.RUnlock()
	m, ok := r.data[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return m, nil
}

// Put adds movie metadata for a given movie id.
func (r *Repository) Put(ctx context.Context, id string, metadata *model.Metadata) error {
	_, span := otel.Tracer(tracerID).Start(ctx, "Repository/Put")
	defer span.End()
	r.Lock()
	defer r.Unlock()
	r.data[id] = metadata
	return nil
}
