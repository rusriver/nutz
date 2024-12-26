package controlflow_test

import (
	"fmt"
	"testing"

	"github.com/rusriver/nutz/controlflow"
)

func Test_03(t *testing.T) {

	fmt.Println("--1")

	controlflow.Try(func() {
		fmt.Println("ok")
	}).Catch(func(e *controlflow.Exception) {
		fmt.Println("exception", e.X)
	})

	fmt.Println("--2")

	controlflow.Try(func() {
		panic("fail")
	}).Catch(func(e *controlflow.Exception) {
		fmt.Println("exception", e.X)
	})

}
