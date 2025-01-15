package clickhouseclient

import (
	"os"
	"strings"
	"unicode"
)

type Program string

func SQL(prgText string, replacerF func(s string) string) *Program {
	prgText = strings.TrimRightFunc(prgText, unicode.IsSpace)
	prgText = os.Expand(prgText, replacerF)
	return (*Program)(&prgText)
}
