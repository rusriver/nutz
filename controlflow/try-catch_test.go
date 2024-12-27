package controlflow_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/rusriver/nutz/controlflow"
)

func Test_03(t *testing.T) {

	fmt.Println("--1")

	controlflow.Try(func() (err error) {
		fmt.Println("ok")
		return
	}).Catch(func(e *controlflow.Exception) {
		fmt.Println("exception", e.X)
	})

	fmt.Println("--2")

	controlflow.Try(func() (err error) {
		panic("fail")
	}).Catch(func(e *controlflow.Exception) {
		fmt.Println("exception", e.X)
	})

	fmt.Println("--3")

	controlflow.Try(func() (err error) {
		return errors.New("error 3")
	}).Catch(func(e *controlflow.Exception) {
		fmt.Println("exception", e.X)
	})

}
