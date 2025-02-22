package factory

import (
	"fmt"

	"github.com/rusriver/nutz/controlflow"
)

// This has absolutely nothing to do with the bullshit "factory pattern" popularized by
// the book "Design Patterns: Elements of Reusable Object-Oriented Software", by
// Gamma, Erich; Helm, Richard; Johnson, Ralph; Vlissides, John (1995), Addison-Wesley,
// ISBN 0-201-63361-2. Nor is it the https://www.google.com/search?q=oop+factory+pattern.
// It has nothing to do with Gang of Four, or Gang of Five, or any other Gang of Bullshit.
// Instead, it specifically implements this:
// https://www.google.com/search?q=factory+production+line&udm=2
// https://en.wikipedia.org/wiki/Production_line
// https://en.wikipedia.org/wiki/Assembly_line
// The Factory is thread-safe for operation, but thread-unsafe for modification.
type Factory[PRODUCT any] struct {
	productionLines   map[string]*Line[PRODUCT]
	Hook_PreLine      func(pprec *ProductProgressRecord[PRODUCT])
	Hook_AfterLine    func(pprec *ProductProgressRecord[PRODUCT])
	Hook_PreStation   func(pprec *ProductProgressRecord[PRODUCT])
	Hook_AfterStation func(pprec *ProductProgressRecord[PRODUCT])
}

func NewFactory[PRODUCT any]() Factory[PRODUCT] {
	f := Factory[PRODUCT]{
		productionLines: map[string]*Line[PRODUCT]{},
	}
	return f
}

func (f Factory[PRODUCT]) AddLineOfObjects(name string) *Line[PRODUCT] {
	line := newLineOfObjects[PRODUCT]()
	f.productionLines[name] = line
	return line
}

func (f Factory[PRODUCT]) AddLineOfFuncs(name string) *Line[PRODUCT] {
	line := newLineOfFuncs[PRODUCT]()
	f.productionLines[name] = line
	return line
}

func (f Factory[PRODUCT]) Line(name string) (line *Line[PRODUCT]) {
	if line, ok := f.productionLines[name]; ok {
		return line
	} else {
		controlflow.Raise1(fmt.Sprintf("55ade7fa04a4ec89729a line '%v' missing", name))
	}
	return
}

// To be used in a user-code managed loop, e.g. when each per-line progress
// gets persistized. Alternatively, consider using ProcessFull().
func (f Factory[PRODUCT]) ProcessOneLine(pprec *ProductProgressRecord[PRODUCT]) {
	if len(pprec.JumpToLine) > 0 {
		pprec.CurrentPosition.Line = pprec.JumpToLine
		pprec.CurrentPosition.Station = 0
		pprec.JumpToLine = ""
	}
	line := f.Line(pprec.CurrentPosition.Line)

	if f.Hook_PreLine != nil {
		f.Hook_PreLine(pprec)
	}
	defer func() {
		if f.Hook_AfterLine != nil {
			f.Hook_AfterLine(pprec)
		}
	}()
	switch line.typeSelector {
	case LineType_Func:
		realLine := *line.line_Funcs
		for i := 0; i < len(realLine); i++ {
			pprec.CurrentPosition.Station = i
			if f.Hook_PreStation != nil {
				f.Hook_PreStation(pprec)
			}
			realLine[i](pprec)
			if f.Hook_AfterStation != nil {
				f.Hook_AfterStation(pprec)
			}
			if pprec.FullStopNow {
				return
			}
			if len(pprec.JumpToLine) > 0 {
				return
			}
		}
	case LineType_Object:
		realLine := *line.line_Objects
		for i := 0; i < len(realLine); i++ {
			pprec.CurrentPosition.Station = i
			if f.Hook_PreStation != nil {
				f.Hook_PreStation(pprec)
			}
			realLine[i].Do(pprec)
			if f.Hook_AfterStation != nil {
				f.Hook_AfterStation(pprec)
			}
			if pprec.FullStopNow {
				return
			}
			if len(pprec.JumpToLine) > 0 {
				return
			}
		}
	}
	pprec.FullStopNow = true
}

func (f Factory[PRODUCT]) ProcessFull(pprec *ProductProgressRecord[PRODUCT]) {
	for {
		f.ProcessOneLine(pprec)
		if pprec.FullStopNow {
			break
		}
	}
}
