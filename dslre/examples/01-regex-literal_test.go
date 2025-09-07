package examples_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/rusriver/nutz/controlflow"
	"github.com/rusriver/nutz/dslre"
	"github.com/rusriver/nutz/dslre/simplecache"
	loggeterrzl "github.com/rusriver/nutz/loggeterr-zl"
)

func Test_RegexLiteral_001(t *testing.T) {

	tracer := loggeterrzl.New(nil)

	// init the DSL registry
	dslreRegex := dslre.NewRegistry(tracer, func(r *dslre.Registry[*regexp.Regexp, any, *regexp.Regexp]) {
		r.InternalCache = simplecache.NewCache[*regexp.Regexp]()
		r.InternalParserFunc = func(literal string) (pp *regexp.Regexp) {
			fmt.Println("++21 compile")
			pp = controlflow.Must(regexp.Compile(literal))
			return
		}
		r.InternalPackerFunc0 = func(pp *regexp.Regexp) (res *regexp.Regexp) {
			return pp
		}
	})

	// the actual program
	{
		res := dslreRegex.ParsegetL(`\sa\s`).MatchString("this is a test")
		fmt.Println(res)
	}
	{
		res := dslreRegex.ParsegetL(`\sa\s`).MatchString("this is a test")
		fmt.Println(res)
	}
	{
		res := dslreRegex.ParsegetL(`\sa\s`).MatchString("this is a test")
		fmt.Println(res)
	}
	{
		res := dslreRegex.ParsegetL(`\sa\s`).MatchString("this is a test")
		fmt.Println(res)
	}
	{
		res := dslreRegex.ParsegetL(`\sa\s`).MatchString("this is a test")
		fmt.Println(res)
	}
}
