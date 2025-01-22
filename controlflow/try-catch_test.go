package controlflow_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

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
		controlflow.Raise("fail")
		return
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

func Test_03_retry(t *testing.T) {

	fmt.Println("--1")
	i := 1
	controlflow.Retry(func() (err error) {
		fmt.Println("some work", i)
		if i >= 5 {
			return nil
		}
		i++
		controlflow.Raise("fail")
		return
	}, func(e *controlflow.Exception) (end bool) {
		fmt.Printf("attempt %v failed: %v; retry\n", e.AttemptNb, e.Err)
		time.Sleep(time.Millisecond * 500)
		return
	}).Catch(func(e *controlflow.Exception) {
		fmt.Println("1 FAIL;", e.Err)
	}).OK(func(e *controlflow.Exception) {
		fmt.Println("1 OK")
	})

	fmt.Println("--2")
	i = 1
	controlflow.Retry(func() (err error) {
		fmt.Println("some work", i)
		if i >= 5 {
			return nil
		}
		i++
		controlflow.Raise("fail")
		return
	}, func(e *controlflow.Exception) (end bool) {
		if e.AttemptNb < 3 {
			fmt.Printf("attempt %v failed: %v; retry\n", e.AttemptNb, e.Err)
			time.Sleep(time.Millisecond * 500)
			return
		} else {
			fmt.Printf("attempt %v failed: %v; yield\n", e.AttemptNb, e.Err)
			end = true
			return
		}
	}).Catch(func(e *controlflow.Exception) {
		fmt.Println("2 FAIL;", e.Err)
	}).OK(func(e *controlflow.Exception) {
		fmt.Println("2 OK")
	})
}
