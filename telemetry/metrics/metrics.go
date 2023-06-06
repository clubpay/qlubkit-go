package qmetrics

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

type Metric struct {
	shutdownFunc func(ctx context.Context) error
	envGauges    []string
}

func New(opts ...Option) (*Metric, error) {
	m := Metric{}

	if len(opts) == 0 {
		return &m, nil
	}

	for _, opt := range opts {
		err := opt(&m)
		if err != nil {
			return nil, err
		}
	}

	_, err := Meter("qmetrics").
		Int64ObservableGauge("env_gauge",
			metric.WithInt64Callback(
				func(ctx context.Context, observer metric.Int64Observer) error {
					attrs := make([]attribute.KeyValue, 0, len(m.envGauges))
					for _, env := range m.envGauges {
						attrs = append(attrs, attribute.String(env, os.Getenv(env)))
					}
					observer.Observe(int64(len(m.envGauges)), metric.WithAttributes(attrs...))

					return nil
				},
			),
		)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (m *Metric) prometheusExporter(port int) error {
	registry := prom.NewRegistry()
	registry.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	exp, err := prometheus.New(
		prometheus.WithRegisterer(registry),
	)
	if err != nil {
		return err
	}
	otel.SetMeterProvider(
		sdkmetric.NewMeterProvider(sdkmetric.WithReader(exp)),
	)

	go servePrometheus(registry, port)

	return nil
}

func (m *Metric) stdExporter() error {
	enc := json.NewEncoder(os.Stdout)
	exp, err := stdoutmetric.New(stdoutmetric.WithEncoder(enc))
	if err != nil {
		return fmt.Errorf("creating stdoutmetric exporter: %w", err)
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(exp),
		),
	)
	m.shutdownFunc = mp.Shutdown
	otel.SetMeterProvider(mp)

	return nil
}

func (m *Metric) Shutdown(ctx context.Context) error {
	if m.shutdownFunc != nil {
		return m.shutdownFunc(ctx)
	}

	return nil
}
