package logger_test

import (
	"fmt"
	"testing"

	"github.com/rusriver/nutz/logger"
)

func Test_Msgtag_01(t *testing.T) {
	s := struct {
		Msgtag *logger.Msgtag
	}{}

	m := s.Msgtag.With("4174e39370e4")

	fmt.Println(m.String())
}
