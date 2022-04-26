package memory

import (
	"context"

	"movieexample.com/rating/pkg/model"
)

// Repository defines a memory movie matadata repository.
type Repository struct {
	data map[model.RecordType]map[model.RecordID]map[model.UserID]model.RatingValue
}

// New creates a new memory repository.
func New() *Repository {
	return &Repository{map[model.RecordType]map[model.RecordID]map[model.UserID]model.RatingValue{}}
}

// Get retrieves movie metadata for by movie id.
func (r *Repository) GetAll(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]*model.Rating, error) {
	if _, ok := r.data[recordType]; !ok {
		return []*model.Rating{}, nil
	}
	if _, ok := r.data[recordType][recordID]; !ok {
		return []*model.Rating{}, nil
	}
	var res []*model.Rating
	for userID, value := range r.data[recordType][recordID] {
		res = append(res, &model.Rating{UserID: userID, Value: value})
	}
	return res, nil
}

// Put adds a rating for a given record.
func (r *Repository) Put(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	if _, ok := r.data[recordType]; !ok {
		r.data[recordType] = map[model.RecordID]map[model.UserID]model.RatingValue{}
	}
	if _, ok := r.data[recordType][recordID]; !ok {
		r.data[recordType][recordID] = map[model.UserID]model.RatingValue{}
	}
	r.data[recordType][recordID][rating.UserID] = rating.Value
	return nil
}
