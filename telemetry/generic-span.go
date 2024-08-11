package telemetry

import "time"

type SpanCreator[T any] struct {
	StarterFunc  func(s *Span[T])
	FinisherFunc func(s *Span[T])
}

func NewSpanCreator[T any](optFuncs ...func(c *SpanCreator[T])) (c *SpanCreator[T]) {
	c = &SpanCreator[T]{}
	for _, of := range optFuncs {
		of(c)
	}
	return
}

type Span[T any] struct {
	starterFunc  func(s *Span[T])
	finisherFunc func(s *Span[T])
	T0, T1       time.Time // the time bounds of the span
	UserAttrs    T
}

func (c *SpanCreator[T]) NewSpan(of ...func(s *Span[T])) (s *Span[T]) {
	s = &Span[T]{
		starterFunc:  c.StarterFunc,
		finisherFunc: c.FinisherFunc,
	}
	for _, f := range of {
		f(s)
	}
	return
}

func (s *Span[T]) Start() {
	if s.starterFunc != nil {
		s.starterFunc(s)
	}
}

// only calls a finisher, if T0 was set, else ignores it
func (s *Span[T]) Finish(of ...func(s *Span[T])) {
	for _, f := range of {
		f(s)
	}
	if s.finisherFunc != nil {
		if !s.T0.IsZero() {
			s.finisherFunc(s)
		}
		s.finisherFunc = nil // make sure the finisher is called only once
	}
}
