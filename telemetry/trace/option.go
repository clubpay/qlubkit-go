package qtrace

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

func WithTerminal(pretty bool) Option {
	return func(t *Tracer) {
		if pretty {
			t.exp = expSTDPretty
		} else {
			t.exp = expSTD
		}
	}
}
