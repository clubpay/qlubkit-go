package qtrace

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"go.opentelemetry.io/otel/trace"
)

type exporter string

const (
	expOTLP      exporter = "otlp"
	expSTD       exporter = "std"
	expSTDPretty exporter = "std-pretty"
)

type Tracer struct {
	exp       exporter
	endpoint  string
	tp        trace.TracerProvider
	tpWrapper func(trace.TracerProvider) trace.TracerProvider
	sampler   sdktrace.Sampler
}

func New(serviceName string, opts ...Option) (*Tracer, error) {
	t := Tracer{
		sampler: sdktrace.ParentBased(defaultSampler),
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

	if t.tpWrapper != nil {
		otel.SetTracerProvider(t.tpWrapper(t.tp))
	} else {
		otel.SetTracerProvider(t.tp)
	}

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

func (t Tracer) Shutdown(ctx context.Context) error {
	if t.tp == nil {
		return nil
	}

	if v, ok := t.tp.(interface{ Shutdown(context.Context) error }); ok {
		return v.Shutdown(ctx)
	}
	
	return nil
}

var defaultSampler = NewSampler("QlubSampler").AddDrop(
	"GET /health-check",
	"GET /health-check-external",
	"GET /heartbeat/:country/:workerID",
	"PaymentSubscribe",
	"Ping",
	"GET /heartbeat/{country}/{workerID}",
)

type CustomSampler struct {
	desc        string
	dropsByName map[string]struct{}
}

func NewSampler(desc string) *CustomSampler {
	return &CustomSampler{
		desc:        desc,
		dropsByName: make(map[string]struct{}),
	}
}

func (t *CustomSampler) AddDrop(name ...string) *CustomSampler {
	for _, n := range name {
		t.dropsByName[n] = struct{}{}
	}

	return t
}

func (t *CustomSampler) ShouldSample(p sdktrace.SamplingParameters) sdktrace.SamplingResult {
	_, ok := t.dropsByName[p.Name]
	if ok {
		return sdktrace.SamplingResult{
			Decision:   sdktrace.Drop,
			Tracestate: trace.SpanContextFromContext(p.ParentContext).TraceState(),
		}
	}

	return sdktrace.SamplingResult{
		Decision:   sdktrace.RecordAndSample,
		Tracestate: trace.SpanContextFromContext(p.ParentContext).TraceState(),
	}
}

func (t *CustomSampler) Description() string {
	return t.desc
}
