package factory_test

import (
	"testing"

	"github.com/rusriver/nutz/controlflow/factory"
)

type MyCtx1 struct{}

func station1(pprec *factory.ProductProgressRecord[*MyCtx1]) {}

func Test_Factory_01(t *testing.T) {
	fac := factory.NewFactory[*MyCtx1]()
	fac.AddLineOfObjects("01")

	// This is an example how you can convert station func to station object, and then
	// add it as an object. Though this is interesting example, there's better way using
	// the "Funcs" type of line, and add function stations directly.

	fac.Line("01").AddStationAsObject(factory.StationFunc[*MyCtx1](station1))
}
