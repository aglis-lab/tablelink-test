package tracer

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func RecordError(span trace.Span, err error) {
	span.SetAttributes(attribute.Bool("error", true))
	span.RecordError(err)
}
