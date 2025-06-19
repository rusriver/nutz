package controlflow_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/rusriver/nutz/controlflow"
	"github.com/rusriver/nutz/loggeterr-api/trace"
)

func Test_03(t *testing.T) {

	fmt.Println("--1")

	controlflow.Try(func() (err any) {
		fmt.Println("ok")
		return
	}).Catch(func(e *controlflow.Event) {
		fmt.Println("exception", e.String())
	})

	fmt.Println("--2")

	controlflow.Try(func() (err any) {
		controlflow.RaiseCause(nil, trace.BacktracePointf("27", "fail"))
		return
	}).Catch(func(e *controlflow.Event) {
		fmt.Println("exception", e.String())
	})

	fmt.Println("--3")

	controlflow.Try(func() (err any) {
		return errors.New("error 3")
	}).Catch(func(e *controlflow.Event) {
		fmt.Println("exception", e.String())
	})

}

func Test_03_retry(t *testing.T) {

	fmt.Println("--1")
	i := 1
	controlflow.Retry(func() (err any) {
		fmt.Println("some work", i)
		if i >= 5 {
			return nil
		}
		i++
		controlflow.RaiseCause(nil, trace.BacktracePointf("53", "fail"))
		return
	}, func(e *controlflow.Event) (end bool) {
		fmt.Printf("attempt %v failed: %v; retry\n", e.RetryAttemptNb, e)
		time.Sleep(time.Millisecond * 500)
		return
	}).Catch(func(e *controlflow.Event) {
		fmt.Println("1 FAIL;", e)
	}).OK(func() {
		fmt.Println("1 OK")
	})

	fmt.Println("--2")
	i = 1
	controlflow.Retry(func() (err any) {
		fmt.Println("some work", i)
		if i >= 5 {
			return nil
		}
		i++
		controlflow.RaiseCause(nil, trace.BacktracePointf("73", "fail"))
		return
	}, func(e *controlflow.Event) (end bool) {
		if e.RetryAttemptNb < 3 {
			fmt.Printf("attempt %v failed: %v; retry\n", e.RetryAttemptNb, e)
			time.Sleep(time.Millisecond * 500)
			return
		} else {
			fmt.Printf("attempt %v failed: %v; yield\n", e.RetryAttemptNb, e)
			end = true
			return
		}
	}).Catch(func(e *controlflow.Event) {
		fmt.Println("2 FAIL;", e)
	}).OK(func() {
		fmt.Println("2 OK")
	})
}

func Test_03_raise2(t *testing.T) {
	controlflow.Try(func() (err any) {
		controlflow.Try(func() (err any) {
			controlflow.Try(func() (err any) {
				controlflow.Try(func() (err any) {
					panic("the problem cause")
				}).Catch(func(e *controlflow.Event) {
					e.Reraise(&trace.BacktracePoint{"1e05a52a098f", nil})
				})
				return
			}).Catch(func(e *controlflow.Event) {
				e.Reraise(&trace.BacktracePoint{"71b9f65e73f1", nil})
			})
			return
		}).Catch(func(e *controlflow.Event) {
			e.Reraise(&trace.BacktracePoint{"927136578a67", nil})
		})
		return
	}).Catch(func(e *controlflow.Event) {
		fmt.Println("46e59cc24375", e)
	})
}

func Test_04_raise2(t *testing.T) {
	controlflow.Try(func() (err any) {
		controlflow.Try(func() (err any) {
			controlflow.Try(func() (err any) {
				controlflow.Try(func() (err any) {
					e := errors.New("a17674b6 the problem cause")
					panic(e)
				}).Catch(func(e *controlflow.Event) {
					e.Reraise(&trace.BacktracePoint{"1e05a52a098f", nil})
				})
				return
			}).Catch(func(e *controlflow.Event) {
				e.Reraise(&trace.BacktracePoint{"71b9f65e73f1", nil})
			})
			return
		}).Catch(func(e *controlflow.Event) {
			e.Reraise(&trace.BacktracePoint{"927136578a67", nil})
		})
		return
	}).Catch(func(e *controlflow.Event) {
		fmt.Println("46e59cc24375", e)
	})
}
