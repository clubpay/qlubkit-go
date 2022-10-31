package qmetrics

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric"
)

type Metric struct {
	shutdownFunc func(ctx context.Context) error
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

	return &m, nil
}

func (m *Metric) prometheusExporter(port int) error {
	exp := prometheus.New()
	global.SetMeterProvider(
		metric.NewMeterProvider(metric.WithReader(exp)),
	)

	go servePrometheus(exp.Collector, port)

	return nil
}

func (m *Metric) stdExporter() error {
	enc := json.NewEncoder(os.Stdout)
	exp, err := stdoutmetric.New(stdoutmetric.WithEncoder(enc))
	if err != nil {
		return fmt.Errorf("creating stdoutmetric exporter: %w", err)
	}
	mp := metric.NewMeterProvider(
		metric.WithReader(
			metric.NewPeriodicReader(exp),
		),
	)
	m.shutdownFunc = mp.Shutdown
	global.SetMeterProvider(mp)

	return nil
}

func (m *Metric) Shutdown(ctx context.Context) error {
	if m.shutdownFunc != nil {
		return m.shutdownFunc(ctx)
	}

	return nil
}
