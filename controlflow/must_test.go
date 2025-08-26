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

		values := controlflow.Must(controlflow.Mustvep(
			url.ParseQuery("a=1&s=2;"),
		), ft, "b17b648bbee12b797619")

		values.Add("extra", "5")
		fmt.Println(values.Encode())

		{
			values, err := url.ParseQuery("a=1&s=2;")
			if err != nil {
				controlflow.RaiseCause(ft, &trace.BacktracePoint{"71e733467f4c", err.Error()})
			}

			_ = values
		}
		return
	}).Catch(func(e *controlflow.Event) {
		fmt.Println(e)
	})
}
