package metrics

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
)

type Metric struct {
	mp metric.MeterProvider
}

func NewMetric(opts ...Option) (*Metric, error) {
	m := Metric{}

	for _, opt := range opts {
		err := opt(&m)
		if err != nil {
			return nil, err
		}
	}

	global.SetMeterProvider(m.mp)

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
