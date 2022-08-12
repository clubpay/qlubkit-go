package job

import (
	"context"
	"sync"
)

type Context struct {
	ctx context.Context
	kvl sync.Mutex
	kv  map[string]interface{}
}

func newCtx(ctx context.Context) *Context {
	return &Context{
		ctx: ctx,
		kv:  map[string]interface{}{},
	}
}

func (ctx *Context) Context() context.Context {
	return ctx.ctx
}

func (ctx *Context) Set(key string, value interface{}) {
	ctx.kvl.Lock()
	ctx.kv[key] = value
	ctx.kvl.Unlock()
}

func (ctx *Context) Get(key string) interface{} {
	ctx.kvl.Lock()
	v := ctx.kv[key]
	ctx.kvl.Unlock()

	return v
}

func (ctx *Context) GetString(key string) string {
	return ctx.GetStringOr(key, "")
}

func (ctx *Context) GetStringOr(key string, or string) string {
	v, ok := ctx.Get(key).(string)
	if !ok {
		return or
	}

	return v
}

func (ctx *Context) GetInt64(key string) int64 {
	return ctx.GetInt64Or(key, 0)
}

func (ctx *Context) GetInt64Or(key string, or int64) int64 {
	v, ok := ctx.Get(key).(int64)
	if !ok {
		return or
	}

	return v
}

func (ctx *Context) GetInt32(key string) int32 {
	return ctx.GetInt32Or(key, 0)
}

func (ctx *Context) GetInt32Or(key string, or int32) int32 {
	v, ok := ctx.Get(key).(int32)
	if !ok {
		return or
	}

	return v
}

func (ctx *Context) GetUint64(key string) uint64 {
	return ctx.GetUint64Or(key, 0)
}

func (ctx *Context) GetUint64Or(key string, or uint64) uint64 {
	v, ok := ctx.Get(key).(uint64)
	if !ok {
		return or
	}

	return v
}

func (ctx *Context) GetUint32(key string) uint32 {
	return ctx.GetUint32Or(key, 0)
}

func (ctx *Context) GetUint32Or(key string, or uint32) uint32 {
	v, ok := ctx.Get(key).(uint32)
	if !ok {
		return or
	}

	return v
}
