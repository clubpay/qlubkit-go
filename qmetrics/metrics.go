package qmetrics

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

type Metric struct {
	mp           metric.MeterProvider
	shutdownFunc func(ctx context.Context) error
}

func New(opts ...Option) (*Metric, error) {
	m := Metric{}

	if len(opts) == 0 {
		err := m.stdExporter()
		if err != nil {
			return nil, err
		}

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
	config := prometheus.Config{
		DefaultHistogramBoundaries: []float64{1, 2, 5, 10, 20, 50},
	}
	c := controller.New(
		processor.NewFactory(
			selector.NewWithHistogramDistribution(
				histogram.WithExplicitBoundaries(config.DefaultHistogramBoundaries),
			),
			aggregation.CumulativeTemporalitySelector(),
			processor.WithMemory(true),
		),
	)

	exp, err := prometheus.New(config, c)
	if err != nil {
		return err
	}

	http.HandleFunc("/", exp.ServeHTTP)
	m.mp = exp.MeterProvider()

	go func() {
		_ = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	}()

	return nil
}

func (m *Metric) otlpExporter(ctx context.Context, endPoint string) error {
	exp, err := otlpmetrichttp.New(ctx,
		otlpmetrichttp.WithEndpoint(endPoint),
		otlpmetrichttp.WithInsecure(),
	)
	if err != nil {
		return err
	}

	ctrl := controller.New(
		processor.NewFactory(
			selector.NewWithInexpensiveDistribution(),
			aggregation.CumulativeTemporalitySelector(),
		),
		controller.WithCollectPeriod(time.Second),
		//controller.WithResource(resource.Empty()),
		controller.WithExporter(exp),
	)

	err = ctrl.Start(context.Background())
	if err != nil {
		return err
	}

	m.shutdownFunc = ctrl.Stop
	global.SetMeterProvider(ctrl)

	return nil
}

func (m *Metric) stdExporter() error {
	exporter, err := stdoutmetric.New(stdoutmetric.WithPrettyPrint())
	if err != nil {
		return fmt.Errorf("creating stdoutmetric exporter: %w", err)
	}

	ctrl := controller.New(
		processor.NewFactory(
			selector.NewWithInexpensiveDistribution(),
			exporter,
		),
		controller.WithExporter(exporter),
	)
	if err = ctrl.Start(context.Background()); err != nil {
		log.Fatalf("starting push controller: %v", err)
	}

	m.shutdownFunc = ctrl.Stop
	global.SetMeterProvider(ctrl)

	return nil
}

func (m *Metric) Shutdown(ctx context.Context) error {
	if m.shutdownFunc != nil {
		return m.shutdownFunc(ctx)
	}

	return nil
}
