package job

import "context"

type ErrorHandler interface {
	OnError(ctx context.Context, err error) FailureAction
}

type FailureAction int

const (
	IgnoreAndContinue = 1 << iota
	StopAndExit
	Retry
)

type noRetry struct{}

func (n noRetry) OnError(ctx context.Context, _ error) FailureAction {
	return IgnoreAndContinue
}

type limitRetry struct {
	remaining int
}

func (r *limitRetry) OnError(ctx context.Context, _ error) FailureAction {
	if r.remaining > 0 {
		r.remaining--

		return Retry
	}

	return IgnoreAndContinue
}
