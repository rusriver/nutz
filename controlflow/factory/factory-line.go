package factory

import (
	"github.com/rusriver/nutz/controlflow"
	"github.com/rusriver/nutz/loggeterr-api/trace"
)

const (
	LineType_Undef = iota
	LineType_Object
	LineType_Func
)

type Line[PRODUCT any] struct {
	typeSelector int8
	line_Objects *[]IStation[PRODUCT]
	line_Funcs   *[]StationFunc[PRODUCT]
}

func newLineOfObjects[PRODUCT any]() (L *Line[PRODUCT]) {
	L = &Line[PRODUCT]{}
	L.typeSelector = LineType_Object
	line := []IStation[PRODUCT]{}
	L.line_Objects = &line
	return L
}

func newLineOfFuncs[PRODUCT any]() (L *Line[PRODUCT]) {
	L = &Line[PRODUCT]{}
	L.typeSelector = LineType_Func
	line := []StationFunc[PRODUCT]{}
	L.line_Funcs = &line
	return L
}

func (line *Line[PRODUCT]) AddStationAsObject(s IStation[PRODUCT]) *Line[PRODUCT] {
	switch line.typeSelector {
	case LineType_Object:
		*line.line_Objects = append(*line.line_Objects, s)
	default:
		controlflow.RaiseCause(nil, trace.BacktracePointf("99af009653ac0982d60e", "incompatible station type"))
	}
	return line
}

func (line *Line[PRODUCT]) AddStationAsFunc(s StationFunc[PRODUCT]) *Line[PRODUCT] {
	switch line.typeSelector {
	case LineType_Func:
		*line.line_Funcs = append(*line.line_Funcs, s)
	default:
		controlflow.RaiseCause(nil, trace.BacktracePointf("2830f927ef5e22d213d", "incompatible station type"))
	}
	return line
}

// syntax sugar
func (line *Line[PRODUCT]) EndOfLine() {}

type IStation[PRODUCT any] interface {
	Do(*ProductProgressRecord[PRODUCT])
}

type StationFunc[PRODUCT any] func(*ProductProgressRecord[PRODUCT])

func (sf StationFunc[PRODUCT]) Do(pr *ProductProgressRecord[PRODUCT]) {
	sf(pr)
}

var _ IStation[any] = StationFunc[any](nil)
