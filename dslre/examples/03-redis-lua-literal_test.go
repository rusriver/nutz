package examples_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/rusriver/nutz/controlflow"
	"github.com/rusriver/nutz/dslre"
	"github.com/rusriver/nutz/dslre/simplecache"
	loggeterrzl "github.com/rusriver/nutz/loggeterr-zl"

	"github.com/alicebob/miniredis/v2"
	_ "github.com/proullon/ramsql/driver"
)

func Test_RedisLuaLiteral_001(t *testing.T) {

	tracer := loggeterrzl.New(nil)

	/*
		In this example we demonstrate passing a struct with args into a packer,
		and also the packer acts as an actual runner of the Lua script, i.e. again
		we use the pattern parse-then-pack quite freely, actually as parse-then-exec.
	*/

	rdbServer := miniredis.RunT(t)

	rdb := redis.NewClient(&redis.Options{
		Addr: rdbServer.Addr(),
	})

	type RedisArgs struct {
		Ctx          context.Context
		ScriptRunner redis.Scripter
		Keys         []string
		Args         []any
	}

	// init the DSL registry
	dslreRedis := dslre.NewRegistry(tracer, func(r *dslre.Registry[*redis.Script, *RedisArgs, *redis.Cmd]) {

		r.InternalStorage = Storage
		r.InternalCache = simplecache.NewCache[*redis.Script]()

		r.InternalParserFunc = func(literal string) (pp *redis.Script) {
			fmt.Println("++36 compile")
			pp = redis.NewScript(literal)
			return
		}

		r.InternalPackerFunc1 = func(pp *redis.Script, args *RedisArgs) (res *redis.Cmd) {
			res = pp.Run(args.Ctx, args.ScriptRunner, args.Keys, args.Args...)
			return
		}

	})

	// the actual program
	{
		res := controlflow.Must(
			dslreRedis.ParsegetS1("lit-03.lua",
				&RedisArgs{context.Background(), rdb, []string{"my_counter"}, []any{+1}},
			).Int(),
		)
		fmt.Println(res)
	}
	{
		res := controlflow.Must(
			dslreRedis.ParsegetS1("lit-03.lua",
				&RedisArgs{context.Background(), rdb, []string{"my_counter"}, []any{+1}},
			).Int(),
		)
		fmt.Println(res)
	}
	{
		res := controlflow.Must(
			dslreRedis.ParsegetS1("lit-03.lua",
				&RedisArgs{context.Background(), rdb, []string{"my_counter"}, []any{+1}},
			).Int(),
		)
		fmt.Println(res)
	}
}
