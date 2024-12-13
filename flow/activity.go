package flow

import (
	"context"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"
)

type ActivityFunc[REQ, RES any] func(ctx *ActivityContext[REQ, RES], req REQ) (*RES, error)

type Activity[REQ, RES, InitArg any] struct {
	sdk     *SDK
	Name    string
	Factory func(InitArg) ActivityFunc[REQ, RES]
}

func (a *Activity[REQ, RES, InitArg]) Init(sdk *SDK, initArg InitArg) {
	a.sdk = sdk
	sdk.w.RegisterActivityWithOptions(
		func(ctx context.Context, req REQ) (*RES, error) {
			return a.Factory(initArg)(
				&ActivityContext[REQ, RES]{
					ctx: ctx,
				}, req,
			)
		},
		activity.RegisterOptions{Name: a.Name, SkipInvalidStructFunctions: true},
	)
}

type ExecuteActivityOptions struct {
	ScheduleToCloseTimeout time.Duration
	ScheduleToStartTimeout time.Duration
	StartToCloseTimeout    time.Duration
	RetryPolicy            *RetryPolicy
}

func (a *Activity[REQ, RES, InitArg]) Execute(ctx Context, req REQ, opts ExecuteActivityOptions) Future[RES] {
	if opts.StartToCloseTimeout == 0 {
		opts.StartToCloseTimeout = time.Minute
	}
	if opts.ScheduleToCloseTimeout == 0 {
		opts.ScheduleToCloseTimeout = time.Hour * 24
	}
	ctx = workflow.WithActivityOptions(
		ctx,
		workflow.ActivityOptions{
			TaskQueue:              a.sdk.taskQ,
			ScheduleToCloseTimeout: opts.ScheduleToCloseTimeout,
			ScheduleToStartTimeout: opts.ScheduleToStartTimeout,
			StartToCloseTimeout:    opts.StartToCloseTimeout,
			RetryPolicy:            opts.RetryPolicy,
		},
	)

	return Future[RES]{
		f: workflow.ExecuteActivity(ctx, a.Name, req),
	}
}
