package clickhouseclient_test

import (
	"context"
	"testing"

	"github.com/rusriver/nutz/clickhouseclient"
)

var (
	ENV_DATABASE = ""
	ENV_HOST     = ""
	ENV_PORT     = ""
	ENV_USER     = ""
	ENV_PASSWORD = ""
)

func Test_JustAUsageIdiom(t *testing.T) {
	boot_Context := context.Background()

	creds := func(c *clickhouseclient.Creds) {
		c.Urls = []string{ENV_HOST + ":" + ENV_PORT}
		c.User = ENV_USER
		c.Password = ENV_PASSWORD
		c.Database = ENV_DATABASE
	}

	e := clickhouseclient.Exec(boot_Context, []func(c *clickhouseclient.Creds){creds}, clickhouseclient.SQL(`
		SHOW ${x}
	`, func(s string) string {
		switch s {
		case "x":
			return "TABLES"
		}
		return ""
	}))
	e.Reraise(nil)

}
