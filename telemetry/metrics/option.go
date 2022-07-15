package qmetrics

import "context"

type Option func(t *Metric) error

func WithPrometheus(port int) Option {
	return func(m *Metric) error {
		return m.prometheusExporter(port)
	}
}

func WithOTLP(endPoint string) Option {
	return func(t *Metric) error {
		return t.otlpExporter(context.Background(), endPoint)
	}
}
