package qmetrics

type Option func(t *Metric) error

func WithPrometheus(port int) Option {
	return func(m *Metric) error {
		return m.prometheusExporter(port)
	}
}

func WithStdout() Option {
	return func(t *Metric) error {
		return t.stdExporter()
	}
}
