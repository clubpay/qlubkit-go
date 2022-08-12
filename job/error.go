package job

import "context"

type ErrorHandler interface {
	Retry(ctx context.Context, err error) bool
}

type noRetry struct{}

func (n noRetry) Retry(ctx context.Context, _ error) bool {
	return false
}

type limitRetry struct {
	remaining int
}

func (r *limitRetry) Retry(ctx context.Context, _ error) bool {
	if r.remaining > 0 {
		r.remaining--

		return true
	}

	return false
}
