package controlflow_test

import (
	"context"
	"testing"

	"github.com/rusriver/nutz/controlflow"
	loggeterrapi "github.com/rusriver/nutz/loggeterr-api"
	loggeterrzl "github.com/rusriver/nutz/loggeterr-zl"
)

type TTracerKey struct{}

var TracerKey TTracerKey

func Test_001(t *testing.T) {
	ctx := context.Background()
	tracer := loggeterrzl.New(nil)
	ctxvaluea := &controlflow.CtxValueAccessor[loggeterrapi.ITracer]{Key: TracerKey}

	ctx = ctxvaluea.SetObjectToCtx(ctx, tracer)

	tracer = ctxvaluea.GetObjectFromCtx(ctx)
	_ = tracer
}
