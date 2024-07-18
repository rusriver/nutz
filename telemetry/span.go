package telemetry

import "time"

type HttpSpanCreator struct {
	StarterFunc  func(s *HttpSpan)
	FinisherFunc func(s *HttpSpan, httpCode int, whoFinishes string)
}

func NewHttpSpanCreator(optFuncs ...func(c *HttpSpanCreator)) (c *HttpSpanCreator, err error) {
	c = &HttpSpanCreator{}
	for _, of := range optFuncs {
		of(c)
	}
	return
}

type HttpSpan struct {
	Parent                *HttpSpanCreator
	Id, Php, Path, Method string
	T0, T1                time.Time // when the span was created
}

func (c *HttpSpanCreator) NewSpan(id, php, path, method string) (s *HttpSpan) {
	s = &HttpSpan{
		Parent: c,
		Id:     id,
		Php:    php,
		Path:   path,
		Method: method,
	}
	return
}

func (s *HttpSpan) Start() {
	if s.Parent.StarterFunc != nil {
		s.Parent.StarterFunc(s)
	}
}

// only calls a finisher, if T0 was set, else ignores it
func (s *HttpSpan) Finish(httpCode int, whoFinishes string) {
	if s.Parent.FinisherFunc != nil {
		if !s.T0.IsZero() {
			s.Parent.FinisherFunc(s, httpCode, whoFinishes)
		}
		s.Parent.FinisherFunc = nil // make sure the finisher is called only once
	}
}
