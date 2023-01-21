package otel

import "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"

// WithHoneCombEndpoint sets the Honeycomb endpoint for the OTLP exporter.
func WithHoneCombEndpoint() otlptracegrpc.Option {
	return otlptracegrpc.WithEndpoint("api.honeycomb.io:443")
}

// WithHoneyCombHeader sets the Honeycomb headers for the OTLP exporter.
func WithHoneyCombHeader(honeycombKey string) otlptracegrpc.Option {
	return otlptracegrpc.WithHeaders(map[string]string{
		"x-honeycomb-team": honeycombKey,
	})
}
