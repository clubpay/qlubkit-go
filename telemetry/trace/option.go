package qtrace

import sdktrace "go.opentelemetry.io/otel/sdk/trace"

type Option func(t *Tracer)

func WithOTLP(endPoint string) Option {
	return func(t *Tracer) {
		t.exp = expOTLP
		t.endpoint = endPoint
	}
}

func WithTerminal(pretty bool) Option {
	return func(t *Tracer) {
		if pretty {
			t.exp = expSTDPretty
		} else {
			t.exp = expSTD
		}
	}
}

func WithCustomSampler(s sdktrace.Sampler) Option {
	return func(t *Tracer) {
		t.sampler = s
	}
}
