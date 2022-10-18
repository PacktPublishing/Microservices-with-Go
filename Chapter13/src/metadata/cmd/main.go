package main

import (
	"context"
	"crypto/md5"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/uber-go/tally"
	"github.com/uber-go/tally/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v3"
	"movieexample.com/gen"
	"movieexample.com/metadata/internal/controller/metadata"
	grpchandler "movieexample.com/metadata/internal/handler/grpc"
	"movieexample.com/metadata/internal/repository/memory"
	"movieexample.com/pkg/discovery"
	"movieexample.com/pkg/discovery/consul"
	"movieexample.com/pkg/tracing"
)

const serviceName = "metadata"

func heavyOperation() {
	for {
		token := make([]byte, 1024)
		rand.Read(token)
		md5.New().Write(token)
	}
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	simulateCPULoad := flag.Bool("simulatecpuload", false, "simulate CPU load for profiling")
	flag.Parse()
	if *simulateCPULoad {
		go heavyOperation()
	}

	go func() {
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			logger.Fatal("Failed to start profiler handler", zap.Error(err))
		}
	}()

	f, err := os.Open("base.yaml")
	if err != nil {
		logger.Fatal("Failed to open configuration", zap.Error(err))
	}
	var cfg config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		logger.Fatal("Failed to parse configuration", zap.Error(err))
	}
	port := cfg.API.Port

	logger.Info("Starting the metadata service", zap.Int("port", port))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tp, err := tracing.NewJaegerProvider(cfg.Jaeger.URL, serviceName)
	if err != nil {
		logger.Fatal("Failed to initialize Jaeger provider", zap.Error(err))
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			logger.Fatal("Failed to shut down Jaeger prodiver", zap.Error(err))
		}
	}()
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	reporter := prometheus.NewReporter(prometheus.Options{})
	scope, closer := tally.NewRootScope(tally.ScopeOptions{
		Tags:           map[string]string{"service": "metadata"},
		CachedReporter: reporter,
	}, 10*time.Second)
	defer closer.Close()
	http.Handle("/metrics", reporter.HTTPHandler())
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Prometheus.MetricsPort), nil); err != nil {
			logger.Fatal("Failed to start the metrics handler", zap.Error(err))
		}
	}()

	counter := scope.Tagged(map[string]string{
		"service": "metadata",
	}).Counter("service_started")
	counter.Inc(1)

	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}
	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				logger.Error("Failed to report healthy state", zap.Error(err))
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)
	repo := memory.New()
	ctrl := metadata.New(repo)
	h := grpchandler.New(ctrl)
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))
	if err != nil {
		logger.Fatal("Failed to listen", zap.Error(err))
	}
	srv := grpc.NewServer(grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()))
	reflection.Register(srv)
	gen.RegisterMetadataServiceServer(srv, h)
	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
}
