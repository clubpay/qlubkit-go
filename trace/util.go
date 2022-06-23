package trace

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func Span(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

func NewSpan(instrument, spanName string, ctx context.Context) (context.Context, trace.Span) {
	ctx, span := otel.Tracer(instrument).Start(ctx, spanName)

	return ctx, span
}
