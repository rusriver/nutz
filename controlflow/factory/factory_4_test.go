package factory_test

import (
	"fmt"
	"testing"

	"github.com/rusriver/nutz/controlflow/factory"
	"github.com/stretchr/testify/require"
)

const Test4_Line01 = "LINE-01"
const Test4_Line02 = "LINE-02"
const Test4_Line03 = "LINE-03"
const Test4_Line04 = "LINE-04"

type MyCtx4 struct {
	ID      string
	Text    string
	Count01 int
}

func (ctx *MyCtx4) AsString() string {
	return fmt.Sprintf("%v: %v", ctx.ID, ctx.Text)
}

// This is a more elaborate example of running multiple lines
func Test_Factory_04(t *testing.T) {

	// create the factory
	f := factory.NewFactory[*MyCtx4]()
	f.Hook_PreStation = func(pprec *factory.ProductProgressRecord[*MyCtx4]) {
		fmt.Println("PRE-STATION", pprec.CurrentPosition)
	}
	f.Hook_AfterStation = func(pprec *factory.ProductProgressRecord[*MyCtx4]) {
		fmt.Println("AFTER-STATION", pprec.CurrentPosition)
	}

	f.AddLineOfFuncs(Test4_Line01).
		AddStationAsFunc(test4_line1_station0).
		AddStationAsFunc(test4_line1_station1).
		AddStationAsFunc(test4_line1_station2).
		EndOfLine()

	f.AddLineOfFuncs(Test4_Line02).
		AddStationAsFunc(test4_line2_station0).
		AddStationAsFunc(test4_line2_station1).
		AddStationAsFunc(test4_line2_station2).
		AddStationAsFunc(test4_line2_station3).
		AddStationAsFunc(test4_line2_station4).
		AddStationAsFunc(test4_line2_station5).
		EndOfLine()

	f.AddLineOfFuncs(Test4_Line03).
		AddStationAsFunc(test4_line3_station0).
		EndOfLine()

	f.AddLineOfFuncs(Test4_Line04).
		AddStationAsFunc(test4_line4_station0).
		AddStationAsFunc(test4_line4_station1).
		AddStationAsFunc(test4_line4_station2).
		EndOfLine()

	// complete, let's use it

	var (
		product *MyCtx4
		pprec   *factory.ProductProgressRecord[*MyCtx4]
	)

	product = &MyCtx4{ID: "001"}
	pprec = factory.NewProductProgressRecord(product, Test4_Line01)
	f.ProcessFull(pprec)
	fmt.Println(product.AsString(), pprec.CurrentPosition)
	require.Equal(t, "001: Hello world! [....] Count01 0 [....] Count01 1 [....] Count01 2 [....] Count01 3 [....] Count01 4 [....] Count01 5 --  we did it, congrats!", product.AsString())

	f.Hook_PreStation = nil
	f.Hook_AfterStation = nil

	product = &MyCtx4{ID: "002"}
	pprec = factory.NewProductProgressRecord(product, Test4_Line01)
	f.ProcessFull(pprec)
	fmt.Println(product.AsString(), pprec.CurrentPosition)
	require.Equal(t, "002: Hello world! [....] Count01 0 [....] Count01 1 [....] Count01 2 [....] Count01 3 [....] Count01 4 [....] Count01 5 --  we did it, congrats!", product.AsString())

	product = &MyCtx4{ID: "003"}
	pprec = factory.NewProductProgressRecord(product, Test4_Line02)
	f.ProcessFull(pprec)
	fmt.Println(product.AsString(), pprec.CurrentPosition)
	require.Equal(t, "003:  [....] Count01 0 [....] Count01 1 [....] Count01 2 [....] Count01 3 [....] Count01 4 [....] Count01 5 --  we did it, congrats!", product.AsString())

	product = &MyCtx4{ID: "004"}
	pprec = factory.NewProductProgressRecord(product, Test4_Line03)
	f.ProcessFull(pprec)
	fmt.Println(product.AsString(), pprec.CurrentPosition)
	require.Equal(t, "004:  Count01 0 [....] Count01 1 [....] Count01 2 [....] Count01 3 [....] Count01 4 [....] Count01 5 --  we did it, congrats!", product.AsString())

	product = &MyCtx4{ID: "005"}
	product.Count01 = 2
	pprec = factory.NewProductProgressRecord(product, Test4_Line03)
	f.ProcessFull(pprec)
	fmt.Println(product.AsString(), pprec.CurrentPosition)
	require.Equal(t, "005:  Count01 2 [....] Count01 3 [....] Count01 4 [....] Count01 5 --  we did it, congrats!", product.AsString())

	product = &MyCtx4{ID: "006"}
	pprec = factory.NewProductProgressRecord(product, Test4_Line04)
	f.ProcessFull(pprec)
	fmt.Println(product.AsString(), pprec.CurrentPosition)
	require.Equal(t, "006:  --  we did it, congrats!", product.AsString())

}

func test4_line1_station0(pprec *factory.ProductProgressRecord[*MyCtx4]) {
	pprec.Product.Text += "Hello"
}
func test4_line1_station1(pprec *factory.ProductProgressRecord[*MyCtx4]) {
	pprec.Product.Text += " world"
}
func test4_line1_station2(pprec *factory.ProductProgressRecord[*MyCtx4]) {
	pprec.Product.Text += "!"
	pprec.JumpToLine = Test4_Line02
}

func test4_line2_station0(pprec *factory.ProductProgressRecord[*MyCtx4]) {
	pprec.Product.Text += " ["
}
func test4_line2_station1(pprec *factory.ProductProgressRecord[*MyCtx4]) {
	pprec.Product.Text += "."
}
func test4_line2_station2(pprec *factory.ProductProgressRecord[*MyCtx4]) {
	pprec.Product.Text += "."
}
func test4_line2_station3(pprec *factory.ProductProgressRecord[*MyCtx4]) {
	pprec.Product.Text += "."
}
func test4_line2_station4(pprec *factory.ProductProgressRecord[*MyCtx4]) {
	pprec.Product.Text += "."
}
func test4_line2_station5(pprec *factory.ProductProgressRecord[*MyCtx4]) {
	pprec.Product.Text += "]"
	pprec.JumpToLine = Test4_Line03
}

func test4_line3_station0(pprec *factory.ProductProgressRecord[*MyCtx4]) {
	pprec.Product.Text += fmt.Sprintf(" Count01 %v", pprec.Product.Count01)
	pprec.Product.Count01++
	pprec.JumpToLine = Test4_Line02
	if pprec.Product.Count01 > 5 {
		pprec.JumpToLine = Test4_Line04
	}
}

func test4_line4_station0(pprec *factory.ProductProgressRecord[*MyCtx4]) {
	pprec.Product.Text += " -- "
}
func test4_line4_station1(pprec *factory.ProductProgressRecord[*MyCtx4]) {
	pprec.Product.Text += " we did it, congrats!"
	pprec.FullStopNow = true
}
func test4_line4_station2(pprec *factory.ProductProgressRecord[*MyCtx4]) {
	pprec.Product.Text += "THIS SHOULD NEVER BE PRINTED"
}
