package replicount

import (
	"time"

	"github.com/lithammer/shortuuid/v4"
)

func Id() (id string) {
	id = time.Now().Format("20060102-150405") + "-" + shortuuid.New()
	return
}
