package tracer

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func newGrpcTracerProvider(ctx context.Context, oltpGrpcProvider string, res *resource.Resource) (*trace.TracerProvider, error) {
	// Create Exporter
	conn, err := grpc.DialContext(ctx, oltpGrpcProvider, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	traceExporter, err := newGrpcExporter(ctx, conn)
	if err != nil {
		return nil, err
	}

	// Set up trace provider.
	batchSpanProcessor := trace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := newTraceProvider(res, batchSpanProcessor)

	return tracerProvider, nil
}

func newTraceProvider(res *resource.Resource, bsp trace.SpanProcessor) *trace.TracerProvider {
	return trace.NewTracerProvider(
		trace.WithSampler(trace.ParentBased(trace.AlwaysSample())), // TODO: Change AlwaysSample for performance reason
		trace.WithResource(res),
		trace.WithSpanProcessor(bsp),
	)
}

func newGrpcExporter(ctx context.Context, conn *grpc.ClientConn) (*otlptrace.Exporter, error) {
	return otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
}
