package main

import (
	"context"
	"log"
	"net"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"movieexample.com/gen"
	metadatatest "movieexample.com/metadata/pkg/testutil"
	movietest "movieexample.com/movie/pkg/testutil"
	"movieexample.com/pkg/discovery"
	"movieexample.com/pkg/discovery/memory"
	ratingtest "movieexample.com/rating/pkg/testutil"
)

const (
	metadataServiceName = "metadata"
	ratingServiceName   = "rating"
	movieServiceName    = "movie"

	metadataServiceAddr = "localhost:8081"
	ratingServiceAddr   = "localhost:8082"
	movieServiceAddr    = "localhost:8083"
)

func main() {
	log.Println("Starting the integration test")

	ctx := context.Background()
	registry := memory.NewRegistry()

	log.Println("Setting up service handlers and clients")

	metadataSrv := startMetadataService(ctx, registry)
	defer metadataSrv.GracefulStop()
	ratingSrv := startRatingService(ctx, registry)
	defer ratingSrv.GracefulStop()
	movieSrv := startMovieService(ctx, registry)
	defer movieSrv.GracefulStop()

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	metadataConn, err := grpc.Dial(metadataServiceAddr, opts)
	if err != nil {
		panic(err)
	}
	defer metadataConn.Close()
	metadataClient := gen.NewMetadataServiceClient(metadataConn)

	ratingConn, err := grpc.Dial(ratingServiceAddr, opts)
	if err != nil {
		panic(err)
	}
	defer ratingConn.Close()
	ratingClient := gen.NewRatingServiceClient(ratingConn)

	movieConn, err := grpc.Dial(movieServiceAddr, opts)
	if err != nil {
		panic(err)
	}
	defer movieConn.Close()
	movieClient := gen.NewMovieServiceClient(movieConn)

	log.Println("Saving test metadata via metadata service")

	m := &gen.Metadata{
		Id:          "the-movie",
		Title:       "The Movie",
		Description: "The Movie, the one and only",
		Director:    "Mr. D",
	}

	if _, err := metadataClient.PutMetadata(ctx, &gen.PutMetadataRequest{Metadata: m}); err != nil {
		log.Fatalf("put metadata: %v", err)
	}

	log.Println("Retrieving test metadata via metadata service")

	getMetadataResp, err := metadataClient.GetMetadata(ctx, &gen.GetMetadataRequest{MovieId: m.Id})
	if err != nil {
		log.Fatalf("get metadata: %v", err)
	}
	if diff := cmp.Diff(getMetadataResp.Metadata, m, cmpopts.IgnoreUnexported(gen.Metadata{})); diff != "" {
		log.Fatalf("get metadata after put mismatch: %v", diff)
	}

	log.Println("Getting movie details via movie service")

	wantMovieDetails := &gen.MovieDetails{
		Metadata: m,
	}

	getMovieDetailsResp, err := movieClient.GetMovieDetails(ctx, &gen.GetMovieDetailsRequest{MovieId: m.Id})
	if err != nil {
		log.Fatalf("get movie details: %v", err)
	}
	if diff := cmp.Diff(getMovieDetailsResp.MovieDetails, wantMovieDetails, cmpopts.IgnoreUnexported(gen.MovieDetails{}, gen.Metadata{})); diff != "" {
		log.Fatalf("get movie details after put mismatch: %v", err)
	}

	log.Println("Saving first rating via rating service")

	const userID = "user0"
	const recordTypeMovie = "movie"
	firstRating := int32(5)
	if _, err = ratingClient.PutRating(ctx, &gen.PutRatingRequest{
		UserId:      userID,
		RecordId:    m.Id,
		RecordType:  recordTypeMovie,
		RatingValue: firstRating,
	}); err != nil {
		log.Fatalf("put rating: %v", err)
	}

	log.Println("Retrieving initial aggregated rating via rating service")

	getAggregatedRatingResp, err := ratingClient.GetAggregatedRating(ctx, &gen.GetAggregatedRatingRequest{
		RecordId:   m.Id,
		RecordType: recordTypeMovie,
	})
	if err != nil {
		log.Fatalf("get aggreggated rating: %v", err)
	}

	if got, want := getAggregatedRatingResp.RatingValue, float64(5); got != want {
		log.Fatalf("rating mismatch: got %v want %v", got, want)
	}

	log.Println("Saving second rating via rating service")

	secondRating := int32(1)
	if _, err = ratingClient.PutRating(ctx, &gen.PutRatingRequest{
		UserId:      userID,
		RecordId:    m.Id,
		RecordType:  recordTypeMovie,
		RatingValue: secondRating,
	}); err != nil {
		log.Fatalf("put rating: %v", err)
	}

	log.Println("Saving new aggregated rating via rating service")

	getAggregatedRatingResp, err = ratingClient.GetAggregatedRating(ctx, &gen.GetAggregatedRatingRequest{
		RecordId:   m.Id,
		RecordType: recordTypeMovie,
	})
	if err != nil {
		log.Fatalf("get aggreggated rating: %v", err)
	}

	wantRating := float64((firstRating + secondRating) / 2)
	if got, want := getAggregatedRatingResp.RatingValue, wantRating; got != want {
		log.Fatalf("rating mismatch: got %v want %v", got, want)
	}

	log.Println("Getting updated movie details via movie service")

	getMovieDetailsResp, err = movieClient.GetMovieDetails(ctx, &gen.GetMovieDetailsRequest{MovieId: m.Id})
	if err != nil {
		log.Fatalf("get movie details: %v", err)
	}
	wantMovieDetails.Rating = wantRating
	if diff := cmp.Diff(getMovieDetailsResp.MovieDetails, wantMovieDetails, cmpopts.IgnoreUnexported(gen.MovieDetails{}, gen.Metadata{})); diff != "" {
		log.Fatalf("get movie details after update mismatch: %v", err)
	}

	log.Println("Integration test execution successful")
}

func startMetadataService(ctx context.Context, registry discovery.Registry) *grpc.Server {
	log.Println("Starting metadata service on " + metadataServiceAddr)
	h := metadatatest.NewTestMetadataGRPCServer()
	l, err := net.Listen("tcp", metadataServiceAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	gen.RegisterMetadataServiceServer(srv, h)
	go func() {
		if err := srv.Serve(l); err != nil {
			panic(err)
		}
	}()
	id := discovery.GenerateInstanceID(metadataServiceName)
	if err := registry.Register(ctx, id, metadataServiceName, metadataServiceAddr); err != nil {
		panic(err)
	}
	return srv
}

func startRatingService(ctx context.Context, registry discovery.Registry) *grpc.Server {
	log.Println("Starting rating service on " + ratingServiceAddr)
	h := ratingtest.NewTestRatingGRPCServer()
	l, err := net.Listen("tcp", ratingServiceAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	gen.RegisterRatingServiceServer(srv, h)
	go func() {
		if err := srv.Serve(l); err != nil {
			panic(err)
		}
	}()
	id := discovery.GenerateInstanceID(ratingServiceName)
	if err := registry.Register(ctx, id, ratingServiceName, ratingServiceAddr); err != nil {
		panic(err)
	}
	return srv
}

func startMovieService(ctx context.Context, registry discovery.Registry) *grpc.Server {
	log.Println("Starting movie service on " + movieServiceAddr)
	h := movietest.NewTestMovieGRPCServer(registry)
	l, err := net.Listen("tcp", movieServiceAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	gen.RegisterMovieServiceServer(srv, h)
	go func() {
		if err := srv.Serve(l); err != nil {
			panic(err)
		}
	}()
	id := discovery.GenerateInstanceID(movieServiceName)
	if err := registry.Register(ctx, id, movieServiceName, movieServiceAddr); err != nil {
		panic(err)
	}
	return srv
}
