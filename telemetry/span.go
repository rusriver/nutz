package telemetry

import "time"

type HttpSpanCreator struct {
	ResponseType_0own1proxied bool
	StarterFunc               func(s *HttpSpan)
	FinisherFunc              func(s *HttpSpan)
}

func NewHttpSpanCreator(optFuncs ...func(f *HttpSpanCreator)) (c *HttpSpanCreator, err error) {
	c = &HttpSpanCreator{}
	for _, of := range optFuncs {
		of(c)
	}
	return
}

type HttpSpan struct {
	Parent                *HttpSpanCreator
	Id, Php, Path, Method string
	T0                    time.Time // when the span was created
}

func (c *HttpSpanCreator) NewSpan(id, php, path, method string) (s *HttpSpan) {
	s = &HttpSpan{
		Parent: c,
		Id:     id,
		Php:    php,
		Path:   path,
		Method: method,
		T0:     time.Now(),
	}
	if c.StarterFunc != nil {
		c.StarterFunc(s)
	}
	return
}

func (s *HttpSpan) Finish(httpCode int) {
	if s.Parent.FinisherFunc != nil {
		s.Parent.FinisherFunc(s)
	}
}
