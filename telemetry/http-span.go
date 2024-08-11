package telemetry

import "time"

type HttpSpanCreator struct {
	StarterFunc  func(s *HttpSpan)
	FinisherFunc func(s *HttpSpan, httpCode int, whoFinishes string)
}

// deprecated, please use the generic span instead
func NewHttpSpanCreator(optFuncs ...func(c *HttpSpanCreator)) (c *HttpSpanCreator, err error) {
	c = &HttpSpanCreator{}
	for _, of := range optFuncs {
		of(c)
	}
	return
}

type HttpSpan struct {
	StarterFunc           func(s *HttpSpan)
	FinisherFunc          func(s *HttpSpan, httpCode int, whoFinishes string)
	Id, Php, Path, Method string
	T0, T1                time.Time // the time bounds of the span
	ExtraObject           any       // can be used to attach extra user context, e.g. to report data transferred, other attributes of a span
}

// deprecated, please use the generic span instead
func (c *HttpSpanCreator) NewSpan(id, php, path, method string) (s *HttpSpan) {
	s = &HttpSpan{
		StarterFunc:  c.StarterFunc,
		FinisherFunc: c.FinisherFunc,
		Id:           id,
		Php:          php,
		Path:         path,
		Method:       method,
	}
	return
}

func (s *HttpSpan) Start() {
	if s.StarterFunc != nil {
		s.StarterFunc(s)
	}
}

// only calls a finisher, if T0 was set, else ignores it
func (s *HttpSpan) Finish(httpCode int, whoFinishes string) {
	if s.FinisherFunc != nil {
		if !s.T0.IsZero() {
			s.FinisherFunc(s, httpCode, whoFinishes)
		}
		s.FinisherFunc = nil // make sure the finisher is called only once
	}
}
