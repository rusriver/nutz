package factory_test

import (
	"testing"

	"github.com/rusriver/nutz/controlflow/factory"
)

type MyCtx2 struct{}

func (my *MyCtx2) station1(pprec *factory.ProductProgressRecord[*MyCtx1]) {}
func (my *MyCtx2) station2(pprec *factory.ProductProgressRecord[*MyCtx1]) {}
func (my *MyCtx2) station3(pprec *factory.ProductProgressRecord[*MyCtx1]) {}

func Test_Factory_02(t *testing.T) {
	fac := factory.NewFactory[*MyCtx1]()
	fac.AddLineOfFuncs("01")

	my := &MyCtx2{}

	// This is a demo how to pass methods, bound to object instance, as stations.
	// The base object (*MyCtx) can be used interchengeably with ProgressRecord.Product.

	fac.Line("01").AddStationAsFunc(my.station1)
	fac.Line("01").AddStationAsFunc(my.station2)
	fac.Line("01").AddStationAsFunc(my.station3)
}
