package loggeterrzl

import (
	"encoding/json"
	"regexp"
	"strings"
)

var reStacktraceFormatter01 = regexp.MustCompile(`([^\:])\n(.)`)
var reStacktraceFormatter02 = regexp.MustCompile(`\n\t`)
var reStacktraceFormatter03 = regexp.MustCompile(`\n`)

// Make a one-line representation of Go stracktrace, optimized to be human-readable in one-line form.
// It's useful when reading stacktraces logged into logging systems, where it's represented as a one-liner,
// with line breaks represented as \n, et al.
// Typical usage idiom:
//
//	e.Str("STACKTRACE", logger.StacktraceReadableOneliner(debug.Stack()))
func StacktraceReadableOneliner(stacktrace []byte) string {
	stacktrace = reStacktraceFormatter02.ReplaceAllLiteral(stacktrace, []byte(" at "))
	stacktrace = reStacktraceFormatter01.ReplaceAll(stacktrace, []byte(`$1  | |  $2`))
	stacktrace = reStacktraceFormatter03.ReplaceAllLiteral(stacktrace, []byte("  "))
	return string(stacktrace)
}

func StacktraceAsOnelineJson(stacktrace []byte) []byte {
	lol := [][]string{} // list-of-lists

	lines := strings.Split(string(stacktrace), "\n")
	for i := 1; i < len(lines)-1; i += 2 {
		f := strings.TrimSpace(lines[i])
		w := strings.TrimSpace(lines[i+1])
		lol = append(lol, []string{f, w})
	}

	bb, _ := json.Marshal(lol)
	return bb
}
