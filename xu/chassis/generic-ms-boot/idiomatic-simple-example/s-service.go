package main

import (
	"context"

	loggeterrapi "github.com/rusriver/nutz/loggeterr-api"
	upcommon "github.com/rusriver/nutz/up/up-common"
	genericmsboot "github.com/rusriver/nutz/xu/chassis/generic-ms-boot"
)

type Service struct {
	Name       string
	Tracer     loggeterrapi.ITracer
	Glove      *genericmsboot.Glove
	cancelFunc context.CancelFunc
	//
	SomeUserParam01 string
	SomeUserParam02 int
	SomeUserParam03 bool
}

var _ upcommon.IRemovable = (*Service)(nil)

func NewService(
	name string,
	tracer loggeterrapi.ITracer,
	glove *genericmsboot.Glove,
	of ...func(s *Service),
) *Service {
	s := &Service{
		Name:   name,
		Tracer: tracer,
		Glove:  glove,
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

	if (*glove.DiagnosticFlags)["96bd92c9c7d2"] {
		tracer.InfoEvent("52").SendMsg("DiagnosticFlag 96bd92c9c7d2 activated")
	}

	var ctx context.Context
	ctx, s.cancelFunc = context.WithCancel(context.Background())
	_ = ctx

	return s
}

func (s *Service) GetName() string {
	return s.Name
}
func (s *Service) WindUp(tracer loggeterrapi.ITracer) {
	tracer.InfoEvent("STOP").Send()
	if s.cancelFunc != nil {
		s.cancelFunc()
	}
}
