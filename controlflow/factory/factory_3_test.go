package factory_test

import (
	"fmt"
	"testing"

	"github.com/rusriver/nutz/controlflow/factory"
)

type MyCtx3 struct {
	ID   string
	Text string
}

func (ctx *MyCtx3) AsString() string {
	return fmt.Sprintf("%v: %v", ctx.ID, ctx.Text)
}

// This is a "heelo world" example of using a factory.
func Test_Factory_03(t *testing.T) {
	const MainLine = "main"

	// create the factory
	fac := factory.NewFactory[*MyCtx3]()
	fac.AddLineOfFuncs(MainLine)
	fac.Line(MainLine).AddStationAsFunc(test3_station1)
	fac.Line(MainLine).AddStationAsFunc(test3_station2)
	fac.Line(MainLine).AddStationAsFunc(test3_station3)

	// process an object through it
	product := &MyCtx3{ID: "001"}
	pprec := factory.NewProductProgressRecord(product, MainLine)
	fac.ProcessOneLine(pprec)
	fmt.Println(product.AsString(), pprec.CurrentPosition)

	// process an object through it
	product = &MyCtx3{ID: "002"}
	pprec = factory.NewProductProgressRecord(product, MainLine)
	for {
		fmt.Println(pprec.CurrentPosition)
		fac.ProcessOneLine(pprec)
		fmt.Println("		", pprec.CurrentPosition)
		if pprec.FullStopNow {
			break
		}
	}
	fmt.Println(product.AsString(), pprec.CurrentPosition)

	// process an object through it
	product = &MyCtx3{ID: "003"}
	pprec = factory.NewProductProgressRecord(product, MainLine)
	fac.ProcessFull(pprec)
	fmt.Println(product.AsString(), pprec.CurrentPosition)

}

func test3_station1(pprec *factory.ProductProgressRecord[*MyCtx3]) {
	pprec.Product.Text += "Hello"
}

func test3_station2(pprec *factory.ProductProgressRecord[*MyCtx3]) {
	pprec.Product.Text += " world"
}

func test3_station3(pprec *factory.ProductProgressRecord[*MyCtx3]) {
	pprec.Product.Text += "!"
}
