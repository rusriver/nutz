package loggerapi_test

import (
	"fmt"
	"testing"

	"github.com/rusriver/nutz/logger/loggerapi"
)

func Test_Msgtag_01(t *testing.T) {
	s := struct {
		Msgtag loggerapi.Msgtags
	}{}

	m := s.Msgtag.With("4174e39370e4")

	fmt.Println(m.String())
}
