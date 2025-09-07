package controlflow_test

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/rusriver/nutz/controlflow"
	"github.com/rusriver/nutz/loggeterr-api/trace"
)

func Test_Must_01(t *testing.T) {

	var ft *trace.Forwardtrace

	controlflow.Try(func() (err any) {

		values := controlflow.MustVep(controlflow.VEP(
			url.ParseQuery("a=1&s=2;"),
		), ft, "b17b648bbee12b797619")

		values.Add("extra", "5")
		fmt.Println(values.Encode())

		return
	}).Catch(func(e *controlflow.Event) {
		fmt.Println(e)
	})
}

func Test_Must_02(t *testing.T) {

	var ft *trace.Forwardtrace

	controlflow.Try(func() (er any) {

		values, err := url.ParseQuery("a=1&s=2;")
		if err != nil {
			controlflow.RaiseCause(ft, &trace.BacktracePoint{"71e733467f4c", err.Error()})
		}
		_ = values
		er = err

		return
	}).Catch(func(e *controlflow.Event) {
		fmt.Println(e)
	})
}

func Test_Must_03(t *testing.T) {

	var ft *trace.Forwardtrace

	controlflow.Try(func() (err any) {

		values := controlflow.VEP(url.ParseQuery("a=1&s=2;")).MustWith(ft, "b17b648bbee12b797619")

		values.Add("extra", "5")
		fmt.Println(values.Encode())

		return
	}).Catch(func(e *controlflow.Event) {
		fmt.Println(e)
	})
}

func Test_Must_04(t *testing.T) {

	controlflow.Try(func() (err any) {

		values := controlflow.Must(url.ParseQuery("a=1&s=2;"))
		_ = values
		values = controlflow.Must(url.ParseQuery("a=1&s=22;"))

		values.Add("extra", "5")
		fmt.Println(values.Encode())

		return
	}).Catch(func(e *controlflow.Event) {
		fmt.Println(e)
	})
}
