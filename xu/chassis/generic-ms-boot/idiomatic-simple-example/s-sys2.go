package main

import (
	"context"

	loggeterrapi "github.com/rusriver/nutz/loggeterr-api"
	upcommon "github.com/rusriver/nutz/up/up-common"
	genericmsboot "github.com/rusriver/nutz/xu/chassis/generic-ms-boot"
)

type S_Sys2 struct {
	Name       string
	Tracer     loggeterrapi.ITracer
	Glove      *genericmsboot.Glove
	core       *S_Core
	cancelFunc context.CancelFunc
	//
	SomeUserParam01 string
	SomeUserParam02 int
	SomeUserParam03 bool
}

var _ upcommon.IRemovable = (*S_Sys2)(nil)

func NewSys2(
	name string,
	tracer loggeterrapi.ITracer,
	glove *genericmsboot.Glove,
	of []func(s *S_Sys2),
	core *S_Core,
) *S_Sys2 {
	s := &S_Sys2{
		Name:   name,
		Tracer: tracer,
		Glove:  glove,
		core:   core,
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

	s.Tracer.InfoEvent("9b971a886a1a").SendMsg(
		param1,
		s.SomeUserParam01,
		s.SomeUserParam02,
		s.SomeUserParam03,
	)

	var ctx context.Context
	ctx, s.cancelFunc = context.WithCancel(context.Background())

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
func (s *S_Sys2) GetName() string {
	return s.Name
}
func (s *S_Sys2) WindUp(tracer loggeterrapi.ITracer) {
	tracer.InfoEvent("STOP").Send()
	if s.cancelFunc != nil {
		s.cancelFunc()
	}
}
