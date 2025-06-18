Usage idiom:

	creds := func(c *clickhouseclient.Creds) {
		c.Urls = []string{ENV_HOST + ":" + ENV_PORT}
		c.User = ENV_USER
		c.Password = ENV_PASSWORD
		c.Database = ENV_DATABASE
	}

	err := clickhouseclient.Exec(boot.Context, creds, clickhouseclient.SQL(`
		SHOW ${x}
	`, func(s string) string {
		switch s {
		case "x":
			return "TABLES"
		}
		return ""
	}))
	controlflow.Checkerr(err)

