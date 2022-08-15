package job

import (
	"context"
	"sync/atomic"
)

var nextJobID int64

type Task func(bag *Bag) error

type Job interface {
	ErrorHandler

	ID() int64
	Name() string
	Tasks() []Task
}

func NewJob(name string, opts ...Option) *implJob {
	jb := implJob{
		id:   atomic.AddInt64(&nextJobID, 1),
		name: name,
		eh:   noRetry{},
	}

	for _, o := range opts {
		o(&jb)
	}

	return &jb
}

func (job implJob) AddTask(t ...Task) Job {
	job.tasks = append(job.tasks, t...)

	return job
}

type implJob struct {
	id    int64
	name  string
	tasks []Task
	eh    ErrorHandler
}

var _ Job = (*implJob)(nil)

func (job implJob) ID() int64 {
	return job.id
}

func (job implJob) Name() string {
	return job.name
}

func (job implJob) Tasks() []Task {
	return job.tasks
}

func (job implJob) OnError(ctx context.Context, err error) FailureAction {
	return job.eh.OnError(ctx, err)
}

type Option func(job *implJob)

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
