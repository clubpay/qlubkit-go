package trace

type Option func(t *Tracer)

func WithOTLP(endPoint string) Option {
	return func(t *Tracer) {
		t.exp = expOTLP
		t.endpoint = endPoint
	}
}

func WithJaeger(endPoint string) Option {
	return func(t *Tracer) {
		t.exp = expJaeger
		t.endpoint = endPoint
	}
}
