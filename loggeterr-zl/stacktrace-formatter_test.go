package loggeterrzl_test

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"testing"

	loggeterrzl "github.com/rusriver/nutz/loggeterr-zl"
)

func Test_StacktraceReadableOneliner(t *testing.T) {

	type LL struct {
		Stacktrace string `json:"STACKTRACE"`
	}

	fmt.Println(string(debug.Stack()))
	fmt.Println()

	ll := &LL{
		Stacktrace: string(debug.Stack()),
	}
	bb, _ := json.Marshal(ll)
	fmt.Println(string(bb))
	fmt.Println()

	ll = &LL{
		Stacktrace: loggeterrzl.StacktraceReadableOneliner(debug.Stack()),
	}
	bb, _ = json.Marshal(ll)
	fmt.Println(string(bb))
	fmt.Println()

}

func Test_StacktraceAsOnelineJson_01(t *testing.T) {

	fmt.Println(string(debug.Stack()))
	fmt.Println()

	logger := loggeterrzl.New(nil)

	logger.LogEvent("44").RawStruct("STACKTRACE", loggeterrzl.StacktraceAsOnelineJson(debug.Stack())).Send()

	/*
		{
		    "level": "PURE-EVENT",
		    "STACKTRACE": [
		        [
		            "runtime/debug.Stack()",
		            "/usr/lib/go/src/runtime/debug/stack.go:26 +0x5e"
		        ],
		        [
		            "github.com/rusriver/nutz/loggeterr-zl_test.Test_StacktraceAsOnelineJson_01(0xc0000d6540?)",
		            "/rusriver/rusriver_nutz/loggeterr-zl/stacktrace-formatter_test.go:44 +0xba"
		        ],
		        [
		            "testing.tRunner(0xc0000d6540, 0x5d1468)",
		            "/usr/lib/go/src/testing/testing.go:1792 +0xf4"
		        ],
		        [
		            "created by testing.(*T).Run in goroutine 1",
		            "/usr/lib/go/src/testing/testing.go:1851 +0x413"
		        ]
		    ],
		    "FORWARDTRACE-TAGS": [
		        "44"
		    ],
		    "ETC": [
		        1
		    ],
		    "tsz": "26.16µs",
		    "tsa": "20250806-123406"
		}
	*/
}
