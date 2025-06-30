package controlflow

import (
	"context"

	"github.com/rusriver/nutz/loggeterr-api/trace"
)

type CtxValueAccessor[T any] struct {
	Key                     any
	ConstructIfMissing_Func func() T
}

func (a *CtxValueAccessor[T]) SetObjectToCtx(parent_ctx context.Context, object any) (new_ctx context.Context) {
	new_ctx = context.WithValue(parent_ctx, a.Key, object)
	return
}

func (a *CtxValueAccessor[T]) GetObjectFromCtx(ctx context.Context) (object T) {
	var ok bool
	object, ok = ctx.Value(a.Key).(T)
	if ok {
		return
	} else {
		if a.ConstructIfMissing_Func != nil {
			object = a.ConstructIfMissing_Func()
			return
		} else {
			RaiseCause(nil, &trace.BacktracePoint{"1d0da85e9673", "incorrect use of CtxValueAccessor"})
		}
	}
	return
}
