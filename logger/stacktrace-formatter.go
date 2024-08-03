package logger

import "regexp"

var reStacktraceFormatter01 = regexp.MustCompile(`([^\:])\n(.)`)
var reStacktraceFormatter02 = regexp.MustCompile(`\n\t`)
var reStacktraceFormatter03 = regexp.MustCompile(`\n`)

// Make a one-line representation of Go stracktrace, optimized to be human-readable in one-line form.
// It's useful when reading stacktraces logged into logging systems, where it's represented as a one-liner,
// with line breaks represented as \n, et al.
// Typical usage idiom:
//
//	e.Str("stacktrace", logger.StacktraceReadableOneliner(debug.Stack()))
func StacktraceReadableOneliner(stacktrace []byte) string {
	stacktrace = reStacktraceFormatter02.ReplaceAllLiteral(stacktrace, []byte(" == "))
	stacktrace = reStacktraceFormatter01.ReplaceAll(stacktrace, []byte(`$1  |||  $2`))
	stacktrace = reStacktraceFormatter03.ReplaceAllLiteral(stacktrace, []byte("  "))
	return string(stacktrace)
}
