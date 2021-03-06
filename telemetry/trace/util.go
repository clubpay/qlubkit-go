package qtrace

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func Span(ctx context.Context, attrs ...attribute.KeyValue) trace.Span {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attrs...)

	return span
}

func NewSpan(instrument, spanName string, ctx context.Context) (context.Context, trace.Span) {
	ctx, span := otel.Tracer(instrument).Start(ctx, spanName)

	return ctx, span
}

func Error(span trace.Span, err error) trace.Span {
	if err == nil {
		return span
	}
	span.SetStatus(codes.Error, err.Error())

	return span
}

func Event(span trace.Span, name string, attrs ...attribute.KeyValue) trace.Span {
	span.AddEvent(name, trace.WithAttributes(attrs...))

	return span
}

func EventF(span trace.Span, format string, v ...interface{}) trace.Span {
	span.AddEvent(fmt.Sprintf(format, v...))

	return span
}

var (
	String     = attribute.String
	Int        = attribute.Int
	Int64      = attribute.Int64
	Int64Slice = attribute.Int64Slice
)
