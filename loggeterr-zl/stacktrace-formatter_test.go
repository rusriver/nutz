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
		Stacktrace string
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
