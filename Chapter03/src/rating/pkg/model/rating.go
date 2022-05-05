package model

// RecordID defines a record id.
type RecordID string

// RecordType defines a record type.
type RecordType int

// UserID defines a user id.
type UserID string

// RatingValue defines a value of a rating record.
type RatingValue int

// Rating defines an individual rating created by a user for some record.
type Rating struct {
	UserID UserID      `json:"userId"`
	Value  RatingValue `json:"value"`
}
