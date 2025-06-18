package controlflow

import (
	"fmt"
	"os"

	"github.com/rusriver/nutz/loggeterr-api/trace"
)

func Checkerr2Exit(tag string, err error) {
	if err != nil {
		fmt.Println(tag, err)
		os.Exit(1)
	}
}

// The equivalent of !!! operator in Go2
func Checkerr2Raise(tag string, err error) {
	if err != nil {
		RaiseCause(nil, trace.BacktracePointf(tag, err.Error()))
	}
}
