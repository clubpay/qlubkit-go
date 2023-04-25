package qtrace

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
)

type exporter string

const (
	expOTLP      exporter = "otlp"
	expJaeger    exporter = "jaeger"
	expSTD       exporter = "std"
	expSTDPretty exporter = "std-pretty"
)

type Tracer struct {
	exp      exporter
	endpoint string
	tp       *sdktrace.TracerProvider
	sampler  sdktrace.Sampler
}

func New(serviceName string, opts ...Option) (*Tracer, error) {
	t := Tracer{
		sampler: sdktrace.ParentBased(customSampler{}),
	}

	for _, opt := range opts {
		opt(&t)
	}

	var (
		b   sdktrace.SpanExporter
		err error
	)
	switch t.exp {
	case expOTLP:
		b, err = otlpExporter(t.endpoint)
	case expJaeger:
		b, err = jaegerExporter(t.endpoint)
	case expSTDPretty:
		b, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
	case expSTD:
		b, err = stdouttrace.New()
	default:
		return &t, nil
	}
	if err != nil {
		return nil, err
	}

	t.tp = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(
			b,
			sdktrace.WithMaxQueueSize(4096),
		),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(serviceName),
			),
		),
		sdktrace.WithSampler(t.sampler),
	)

	otel.SetTracerProvider(t.tp)

	return &t, nil
}

func otlpExporter(endPoint string) (sdktrace.SpanExporter, error) {
	c := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(endPoint),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithCompression(otlptracehttp.GzipCompression),
		otlptracehttp.WithRetry(
			otlptracehttp.RetryConfig{
				Enabled:         true,
				InitialInterval: time.Second,
				MaxInterval:     10 * time.Second,
				MaxElapsedTime:  time.Minute,
			},
		),
	)

	return otlptrace.New(context.Background(), c)
}

func jaegerExporter(endPoint string) (sdktrace.SpanExporter, error) {
	return jaeger.New(
		jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint(endPoint),
		),
	)
}

func (t Tracer) Shutdown(ctx context.Context) error {
	if t.tp == nil {
		return nil
	}

	return t.tp.Shutdown(ctx)
}

type customSampler struct{}

func (t customSampler) ShouldSample(p sdktrace.SamplingParameters) sdktrace.SamplingResult {
	switch p.Name {
	case "GET /health-check", "GET /health-check-external",
		"GET /heartbeat/:country/:workerID", "PaymentSubscribe", "Ping":
		return sdktrace.SamplingResult{
			Decision:   sdktrace.Drop,
			Tracestate: trace.SpanContextFromContext(p.ParentContext).TraceState(),
		}
	default:
	}

	return sdktrace.SamplingResult{
		Decision:   sdktrace.RecordAndSample,
		Tracestate: trace.SpanContextFromContext(p.ParentContext).TraceState(),
	}
}

func (t customSampler) Description() string {
	return "ClubPaySampler"
}
