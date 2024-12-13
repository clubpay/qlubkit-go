package flow

import (
	"time"

	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/workflow"
)

type WorkflowContext[REQ, RES any] struct {
	ctx workflow.Context
}

type WorkflowInfo = workflow.Info

func (ctx WorkflowContext[REQ, RES]) Info() *WorkflowInfo {
	return workflow.GetInfo(ctx.ctx)
}

func (ctx WorkflowContext[REQ, RES]) Context() workflow.Context {
	return ctx.ctx
}

func (ctx WorkflowContext[REQ, RES]) Selector(name string) workflow.Selector {
	return workflow.NewNamedSelector(ctx.ctx, name)
}

func (ctx WorkflowContext[REQ, RES]) WaitGroup() workflow.WaitGroup {
	return workflow.NewWaitGroup(ctx.ctx)
}

func (ctx WorkflowContext[REQ, RES]) Log() log.Logger {
	return workflow.GetLogger(ctx.ctx)
}

func (ctx WorkflowContext[REQ, RES]) Sleep(d time.Duration) error {
	return workflow.Sleep(ctx.ctx, d)
}

func SideEffect[REQ, RES, T any](ctx *WorkflowContext[REQ, RES], fn func() T) (T, error) {
	reqEncoded := workflow.SideEffect(
		ctx.Context(),
		func(wctx workflow.Context) any {
			return fn()
		},
	)

	var out T
	err := reqEncoded.Get(&out)

	return out, err
}

func MutableSideEffect[REQ, RES any, T comparable](ctx *WorkflowContext[REQ, RES], id string, fn func() T) (T, error) {
	reqEncoded := workflow.MutableSideEffect(
		ctx.Context(),
		id,
		func(wctx workflow.Context) any {
			return fn()
		},
		func(a, b any) bool {
			return a.(T) == b.(T) //nolint:forcetypeassert
		},
	)

	var out T
	err := reqEncoded.Get(&out)

	return out, err
}
