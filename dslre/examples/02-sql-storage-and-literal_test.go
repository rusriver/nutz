package examples_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/rusriver/nutz/controlflow"
	"github.com/rusriver/nutz/dslre"
	"github.com/rusriver/nutz/dslre/simplecache"
	loggeterrzl "github.com/rusriver/nutz/loggeterr-zl"

	_ "github.com/proullon/ramsql/driver"
)

func Test_SqlStorage_001(t *testing.T) {

	tracer := loggeterrzl.New(nil)

	/*
		This is interesting example of the fact that original concept of parse-then-pack
		can be used for different purposes as well. In this case, the packer doesn't pack
		SQL statement args (in is done by Exec() on the user side), but instead we use
		packer to optionally make a statement of a transaction, if we pass a transaction
		argument.

		This is just an example of the fact that this framework and concept can be used quite
		freely in relation to specific tasks.
	*/

	db := controlflow.VEP(sql.Open("ramsql", "test")).MustWith(tracer.GetForwardtrace(), "a14bcfd3")

	// init the DSL registry
	dslreSql := dslre.NewRegistry(tracer, func(r *dslre.Registry[*sql.Stmt, *sql.Tx, *sql.Stmt]) {
		r.InternalCache = simplecache.NewCache[*sql.Stmt]()
		r.InternalParserFunc = func(literal string) (pp *sql.Stmt) {
			fmt.Println("++34 compile")
			pp = controlflow.VEP(db.Prepare(literal)).MustWith(tracer.GetForwardtrace(), "4f99ee74")
			return
		}
		r.InternalPackerFunc0 = func(pp *sql.Stmt) (res *sql.Stmt) {
			return pp
		}
		r.InternalPackerFunc1 = func(pp *sql.Stmt, tx *sql.Tx) (res *sql.Stmt) {
			return tx.Stmt(pp)
		}
		r.InternalStorage = Storage
	})

	// the actual program
	{
		res := controlflow.Must(
			dslreSql.ParsegetS("lit-01.sql").Exec(),
		)
		fmt.Println("ROWS AFFECTED", controlflow.Must(res.RowsAffected()))
	}
	{
		for i := 1000; i < 1010; i++ {
			res := controlflow.Must(
				dslreSql.ParsegetS("lit-02.sql").Exec(12, i),
			)
			fmt.Println("ROWS AFFECTED", controlflow.Must(res.RowsAffected()))
		}
	}
	{
		for i := 2000; i < 2005; i++ {
			res := controlflow.Must(
				dslreSql.ParsegetS("lit-02.sql").Exec(12, i),
			)
			fmt.Println("ROWS AFFECTED", controlflow.Must(res.RowsAffected()))
		}
	}
	{
		rows := controlflow.Must(
			dslreSql.ParsegetL("select * from user_addresses").Query(),
		)
		for rows.Next() {
			var addressId, userId int
			controlflow.Must(0, rows.Scan(&addressId, &userId))
			fmt.Println(addressId, userId)
		}
		rows.Close()
	}
}
