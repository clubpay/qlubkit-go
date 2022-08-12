package job

import (
	"context"
	"sync/atomic"
)

var nextJobID int64

type Job interface {
	ErrorHandler

	ID() int64
	Func() Func
	// RunAfter returns a list of job ids that needs to be already run before this
	// job can start.
	RunAfter() []int64
}

func NewJob(f Func, opts ...Option) *implJob {
	jb := implJob{
		f:  f,
		id: atomic.AddInt64(&nextJobID, 1),
		eh: noRetry{},
	}

	for _, o := range opts {
		o(&jb)
	}

	return &jb
}

type implJob struct {
	id        int64
	f         Func
	dependsOn []int64
	eh        ErrorHandler
}

var _ Job = (*implJob)(nil)

func (job implJob) ID() int64 {
	return job.id
}

func (job implJob) Func() Func {
	return job.f
}

func (job implJob) RunAfter() []int64 {
	return job.dependsOn
}

func (job implJob) Retry(ctx context.Context, err error) bool {
	return job.eh.Retry(ctx, err)
}

type Option func(job *implJob)

func DependsOn(ids ...int64) Option {
	return func(job *implJob) {
		job.dependsOn = ids
	}
}

func WithMaxRetry(retries int) Option {
	return func(job *implJob) {
		job.eh = &limitRetry{
			remaining: retries,
		}
	}
}

func WithCustomErrorHandler(eh ErrorHandler) Option {
	return func(job *implJob) {
		job.eh = eh
	}
}
