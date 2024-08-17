package telemetry

import (
	"errors"
	"time"
)

type SpanCreator[T any] struct {
	StarterFunc  func(s *Span[T])
	FinisherFunc func(s *Span[T])
	ErrorFunc    func(s *Span[T], err error)
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
	errorFunc    func(s *Span[T], err error)
	T0, T1       time.Time // the time bounds of the span
	UserAttrs    T
}

func (c *SpanCreator[T]) NewSpan(of ...func(s *Span[T])) (s *Span[T]) {
	s = &Span[T]{
		starterFunc:  c.StarterFunc,
		finisherFunc: c.FinisherFunc,
		errorFunc:    c.ErrorFunc,
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
			if !s.T1.After(s.T0) {
				s.T1 = s.T0 // sometimes the span might have the timestamp only, then the duration would be zero
			}
			s.finisherFunc(s)
		} else {
			if s.errorFunc != nil {
				s.errorFunc(s, errors.New("T0 must not be zero"))
			}
		}
		s.finisherFunc = nil // make sure the finisher is called only once
	}
}
