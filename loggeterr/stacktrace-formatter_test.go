package loggeterr_test

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"testing"

	"github.com/rusriver/nutz/loggeterr"
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
		Stacktrace: loggeterr.StacktraceReadableOneliner(debug.Stack()),
	}
	bb, _ = json.Marshal(ll)
	fmt.Println(string(bb))
	fmt.Println()

}
