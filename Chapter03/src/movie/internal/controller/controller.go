package controller

import (
	"context"
	"errors"

	metadatamodel "movieexample.com/metadata/pkg/model"
	"movieexample.com/movie/internal/gateway"
	"movieexample.com/movie/pkg/model"
	ratingmodel "movieexample.com/rating/pkg/model"
)

// ErrNotFound is returned when the movie metadata is not found.
var ErrNotFound = errors.New("not found")

type ratingGateway interface {
	GetAggregatedRating(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType) (float64, error)
	PutRating(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType, rating *ratingmodel.Rating) error
}

type metadataGateway interface {
	Get(ctx context.Context, id string) (*metadatamodel.Metadata, error)
}

// MovieService defines a movie service.
type MovieService struct {
	ratingGateway   ratingGateway
	metadataGateway metadataGateway
}

// New creates a new movie service.
func New(ratingGateway ratingGateway, metadataGateway metadataGateway) *MovieService {
	return &MovieService{ratingGateway, metadataGateway}
}

// Get returns the movie details including the aggregated rating and movie metadata.
func (s *MovieService) Get(ctx context.Context, id string) (*model.MovieDetails, error) {
	metadata, err := s.metadataGateway.Get(ctx, id)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	details := &model.MovieDetails{Metadata: *metadata}
	rating, err := s.ratingGateway.GetAggregatedRating(ctx, ratingmodel.RecordID(id), 1)
	if err != nil && !errors.Is(err, gateway.ErrNotFound) {
		// Just proceed in this case, it's ok not to have ratings yet.
	} else if err != nil {
		return nil, err
	} else {
		details.Rating = &rating
	}
	return details, nil
}
