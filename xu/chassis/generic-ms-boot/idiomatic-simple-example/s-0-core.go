package main

import (
	"context"

	loggeterrapi "github.com/rusriver/nutz/loggeterr-api"
	upcommon "github.com/rusriver/nutz/up/up-common"
	genericmsboot "github.com/rusriver/nutz/xu/chassis/generic-ms-boot"
)

type S_Core struct {
	Name       string
	Tracer     loggeterrapi.ITracer
	CancelFunc context.CancelFunc
	//
	SomeUserParam01 string
	SomeUserParam02 int
	SomeUserParam03 bool
}

var _ upcommon.IRemovable = (*S_Core)(nil)

func NewCore(
	name string,
	tracer loggeterrapi.ITracer,
	glove *genericmsboot.Glove,
	of ...func(s *S_Core),
) *S_Core {
	s := &S_Core{
		Name:   name,
		Tracer: tracer,
	}
	// defaults are here
	s.SomeUserParam01 = "hello"
	s.SomeUserParam02 = 12
	s.SomeUserParam03 = true
	for _, f := range of {
		f(s)
	}

	s.Tracer.InfoEvent("START").Send()

	param1 := glove.ConfDB.P("param-1").String()

	s.Tracer.InfoEvent("0198437cd198").SendMsg(
		param1,
		s.SomeUserParam01,
		s.SomeUserParam02,
		s.SomeUserParam03,
	)

	var ctx context.Context
	ctx, s.CancelFunc = context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ctx.Done():
				return
			}
		}
	}()

	return s
}

func (s *S_Core) GetName() string {
	return s.Name
}
func (s *S_Core) WindUp(tracer loggeterrapi.ITracer) {
	tracer.InfoEvent("STOP").Send()
	if s.CancelFunc != nil {
		s.CancelFunc()
	}
}
