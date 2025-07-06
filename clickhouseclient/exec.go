package clickhouseclient

import (
	"context"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/rusriver/nutz/controlflow"
	"github.com/rusriver/nutz/loggeterr-api/trace"
)

type Creds struct {
	Urls        []string
	Port        string
	Database    string
	User        string
	Password    string
	MaxExecTime time.Duration
}

func Exec(ctx context.Context, creds_func []func(c *Creds), prg *Program) (e *controlflow.Event) {
	c := &Creds{}
	// defaults here
	c.MaxExecTime = time.Minute
	for _, f := range creds_func {
		f(c)
	}

	prgText := string(*prg)
	fmt.Println(prgText)

	dbOpts := &clickhouse.Options{
		Addr: c.Urls,
		Auth: clickhouse.Auth{
			Database: c.Database,
			Username: c.User,
			Password: c.Password,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": int(c.MaxExecTime.Seconds()),
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		DialTimeout:          time.Second * 5,
		ConnOpenStrategy:     clickhouse.ConnOpenRandom,
		BlockBufferSize:      2,
		MaxCompressionBuffer: 500_000,
		ClientInfo: clickhouse.ClientInfo{
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "cli tool"},
			},
		},
	}

	conn, err := clickhouse.Open(dbOpts)
	if err != nil {
		return controlflow.NewEvent(nil, &trace.BacktracePoint{"82b632e1458c", err.Error()})
	}
	v, err := conn.ServerVersion()
	if err != nil {
		return controlflow.NewEvent(nil, &trace.BacktracePoint{"637b217976b1", err.Error()})
	}
	fmt.Println(v, "connected...")
	defer conn.Close()

	stats := &Statistics{}
	controlflow.Try(func() (err any) {
		var el time.Duration
		ctx := clickhouse.Context(ctx, clickhouse.WithProgress(func(p *clickhouse.Progress) {
			if p.Elapsed-el > time.Second*3 {
				el = p.Elapsed
				fmt.Print(".")
			}
			stats.progress.Rows += p.Rows
			stats.progress.Bytes += p.Bytes
			stats.progress.WroteRows += p.WroteRows
			stats.progress.WroteBytes += p.WroteBytes
			stats.progress.Elapsed += p.Elapsed
		}))

		// https://clickhouse.com/docs/en/integrations/go#execution
		fmt.Print("exec..")
		err = conn.Exec(ctx, prgText)
		return
	}).Catch(func(e2 *controlflow.Event) {
		e = e2
	})
	fmt.Println()
	stats.Humanize()
	if e != nil {
		stats.Status = e.String()
	} else {
		stats.Status = "OK"
	}
	fmt.Println(stats.JSON())
	return
}
