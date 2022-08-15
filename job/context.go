package job

import (
	"context"
	"sync"
)

type Bag struct {
	ctx context.Context
	kvl sync.Mutex
	kv  map[string]interface{}
}

func newBag(ctx context.Context) *Bag {
	return &Bag{
		ctx: ctx,
		kv:  map[string]interface{}{},
	}
}

func (ctx *Bag) Context() context.Context {
	return ctx.ctx
}

func (ctx *Bag) Set(key string, value interface{}) {
	ctx.kvl.Lock()
	ctx.kv[key] = value
	ctx.kvl.Unlock()
}

func (ctx *Bag) Get(key string) interface{} {
	ctx.kvl.Lock()
	v := ctx.kv[key]
	ctx.kvl.Unlock()

	return v
}

func (ctx *Bag) GetString(key string) string {
	return ctx.GetStringOr(key, "")
}

func (ctx *Bag) GetStringOr(key string, or string) string {
	v, ok := ctx.Get(key).(string)
	if !ok {
		return or
	}

	return v
}

func (ctx *Bag) GetInt64(key string) int64 {
	return ctx.GetInt64Or(key, 0)
}

func (ctx *Bag) GetInt64Or(key string, or int64) int64 {
	v, ok := ctx.Get(key).(int64)
	if !ok {
		return or
	}

	return v
}

func (ctx *Bag) GetInt32(key string) int32 {
	return ctx.GetInt32Or(key, 0)
}

func (ctx *Bag) GetInt32Or(key string, or int32) int32 {
	v, ok := ctx.Get(key).(int32)
	if !ok {
		return or
	}

	return v
}

func (ctx *Bag) GetUint64(key string) uint64 {
	return ctx.GetUint64Or(key, 0)
}

func (ctx *Bag) GetUint64Or(key string, or uint64) uint64 {
	v, ok := ctx.Get(key).(uint64)
	if !ok {
		return or
	}

	return v
}

func (ctx *Bag) GetUint32(key string) uint32 {
	return ctx.GetUint32Or(key, 0)
}

func (ctx *Bag) GetUint32Or(key string, or uint32) uint32 {
	v, ok := ctx.Get(key).(uint32)
	if !ok {
		return or
	}

	return v
}
