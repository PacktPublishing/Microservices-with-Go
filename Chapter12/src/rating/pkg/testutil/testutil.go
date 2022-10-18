package testutil

import (
	"movieexample.com/gen"
	"movieexample.com/rating/internal/controller/rating"
	grpchandler "movieexample.com/rating/internal/handler/grpc"
	"movieexample.com/rating/internal/repository/memory"
)

// NewTestRatingGRPCServer creates a new rating gRPC server to be used in tests.
func NewTestRatingGRPCServer() gen.RatingServiceServer {
	r := memory.New()
	ctrl := rating.New(r, nil)
	return grpchandler.New(ctrl)
}
