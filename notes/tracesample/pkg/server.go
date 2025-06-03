package pkg

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-logr/zapr"
	// TODO(jeremy): Need to fix this.
	"github.com/jlewi/notes/some/protos"
	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// Server is a grpc server
type Server struct {
	port int

	client v1alpha1.ReconcileServiceClient
	// This is to prevent leaking between requests.
	v1alpha1.UnimplementedLanguageServiceServer
	v1alpha1.UnimplementedReconcileServiceServer
}

// RunServer runs the server on the given port
func RunServer(port int) error {
	log := zapr.NewLogger(zap.L())
	server := &Server{
		port: port,
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)

	v1alpha1.RegisterLanguageServiceServer(grpcServer, server)
	v1alpha1.RegisterReconcileServiceServer(grpcServer, server)

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	// Register health check.
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return errors.Wrapf(err, "failed to listen: %v", err)
	}
	trapInterrupt(grpcServer)

	log.Info("Starting grpc service", "port", port)
	// This will block
	if err := grpcServer.Serve(lis); err != nil {
		return errors.Wrapf(err, "Failed to start the server")
	}

	return nil
}

// trapInterrupt waits for a shutdown signal and shutsdown the server
func trapInterrupt(s *grpc.Server) {
	log := zapr.NewLogger(zap.L())
	sigs := make(chan os.Signal, 10)
	// SIGSTOP and SIGTERM can't be caught; however SIGINT works as expected when using ctl-z
	// to interrupt the process
	signal.Notify(sigs, syscall.SIGINT)

	go func() {
		msg := <-sigs
		log.Info("Received shutdown signal; shutting down gRPC server", "sig", msg)
		s.GracefulStop()
		log.Info("gRPC server shutdown complete")
	}()
}

func (s *Server) CompileAndRun(ctx context.Context, req *v1alpha1.CompileAndRunRequest) (*v1alpha1.CompileAndRunResponse, error) {
	// Get the current span from the context.
	span := trace.SpanFromContext(ctx)
	log := zapr.NewLogger(zap.L())

	log.Info("Received request", "req", req, "spanId", span.SpanContext().SpanID(), "traceId", span.SpanContext().TraceID())
	span.SetAttributes(attribute.String("myattr", "customfoo"))

	// Create the client
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()))

	if err != nil {
		return nil, errors.Wrapf(err, "Failed to connect to server")
	}
	defer conn.Close()

	client := v1alpha1.NewReconcileServiceClient(conn)

	reconcileReq := &v1alpha1.ReconcileRequest{}
	resp, err := client.Reconcile(ctx, reconcileReq)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to reconcile")
	}

	log.Info("Calling http server")

	// Override the defaultclient used by otelhttp so that we use a transport that will report metrics
	otelhttp.DefaultClient = &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	// Get is a helper method that injects a context into the Get request and adds a span around the
	// request
	httpResp, err := otelhttp.Get(ctx, "http://localhost:8080/healthz")
	if err != nil {
		return nil, errors.Wrapf(err, "Http Get request failed")
	}

	// We need to close the body in order to report the trace.
	httpResp.Body.Close()
	log.Info("Http Get request succeeded", "status", httpResp.Status)

	return &v1alpha1.CompileAndRunResponse{
		ExitCode: 1,
		StdOut:   "CompileAndRunCalled: + " + resp.String(),
		StdErr:   "This is stderr",
	}, nil
}

func (s *Server) Reconcile(ctx context.Context, request *v1alpha1.ReconcileRequest) (*v1alpha1.ReconcileResponse, error) {
	return &v1alpha1.ReconcileResponse{
		Files: []*v1alpha1.File{
			{
				Path: "foo.py",
				Contents: &v1alpha1.File_Body{
					Body: "print('Hello World')",
				},
			},
		},
	}, nil
}
