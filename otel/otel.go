package otel

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"google.golang.org/grpc/credentials"
)

// SetupTracing configures the global OpenTelemetry SDK to send trace data.
func SetupTracing(
	ctx context.Context,
	rcs *resource.Resource,
	opts ...otlptracegrpc.Option,
) (shutdown func(), err error) {
	// Configure a new exporter using environment variables for sending data to Honeycomb over gRPC.
	exp, err := newExporter(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize exporter: %w", err)
	}

	// Create a new tracer provider with a batch span processor and the otlp exporter.
	tp := newTraceProvider(exp, rcs)

	// Set the Tracer Provider and the W3C Trace Context propagator as globals
	otel.SetTracerProvider(tp)

	// Register the trace context and baggage propagators so data is propagated across services/processes.
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return func() { _ = tp.Shutdown(ctx) }, nil
}

func newExporter(ctx context.Context, opts ...otlptracegrpc.Option) (*otlptrace.Exporter, error) {
	opts = append(opts,
		otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")),
	)

	client := otlptracegrpc.NewClient(opts...)
	return otlptrace.New(ctx, client)
}

func newTraceProvider(exp *otlptrace.Exporter, rcs *resource.Resource) *sdktrace.TracerProvider {
	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(rcs),
	)
}

// UnspecifiedValue is the value used for unspecified attributes.
const UnspecifiedValue = "unspecified"

// NewResource create a new resource with the default attributes configured from the environment.
func NewResource(attrs ...attribute.KeyValue) *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		append([]attribute.KeyValue{
			semconv.ServiceNameKey.String(tryEnvs("OTEL_SERVICE_NAME", "POD_NAMESPACE")),
			semconv.ServiceNamespaceKey.String(tryEnvs("OTEL_SERVICE_NAMESPACE", "POD_NAMESPACE")),
			semconv.ServiceInstanceIDKey.String(tryEnvs("OTEL_SERVICE_INSTANCE_ID", "POD_NAME")),
			semconv.ServiceVersionKey.String(tryEnvs("OTEL_SERVICE_VERSION")),
			attribute.String("environment", tryEnvs("OTEL_ENVIRONMENT")),
		}, attrs...)...,
	)
}

func tryEnvs(tryEnvs ...string) string {
	for _, env := range tryEnvs {
		if v := os.Getenv(env); v != "" {
			return v
		}
	}
	return UnspecifiedValue
}
