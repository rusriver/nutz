package controlflow

import (
	"context"
)

type CtxValueAccessor[T any] struct {
	Key any
}

func (a *CtxValueAccessor[T]) SetObjectToCtx(parent_ctx context.Context, object any) (new_ctx context.Context) {
	new_ctx = context.WithValue(parent_ctx, a.Key, object)
	return
}

func (a *CtxValueAccessor[T]) GetObjectFromCtx(ctx context.Context) (object T) {
	object = ctx.Value(a.Key).(T) // unsafe intentionally
	return
}
