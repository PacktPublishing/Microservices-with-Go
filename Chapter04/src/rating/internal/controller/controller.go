package controller

import (
	"context"
	"errors"

	"movieexample.com/rating/pkg/model"
)

// ErrNotFound is returned when no ratings are found for a record.
var ErrNotFound = errors.New("not found")

type ratingRepository interface {
	GetAll(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]*model.Rating, error)
	Put(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error
}

// RatingService encapsulates the rating service business logic.
type RatingService struct {
	repo ratingRepository
}

// New creates a rating service.
func New(repo ratingRepository) *RatingService {
	return &RatingService{repo}
}

// GetAggregatedRating returns the aggregated rating for a record or ErrNotFound if there are no ratings for it.
func (s *RatingService) GetAggregatedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
	ratings, err := s.repo.GetAll(ctx, recordID, recordType)
	if err != nil {
		return 0, err
	}
	if len(ratings) == 0 {
		return 0, ErrNotFound
	}
	sum := float64(0)
	for _, r := range ratings {
		sum += float64(r.Value)
	}
	return sum / float64(len(ratings)), nil
}

// PutRating writes a rating for a given record.
func (s *RatingService) PutRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	return s.repo.Put(ctx, recordID, recordType, rating)
}
