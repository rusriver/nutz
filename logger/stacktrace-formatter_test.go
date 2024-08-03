package logger_test

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"testing"

	"github.com/rusriver/nutz/logger"
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
		Stacktrace: logger.StacktraceReadableOneliner(debug.Stack()),
	}
	bb, _ = json.Marshal(ll)
	fmt.Println(string(bb))
	fmt.Println()

}
