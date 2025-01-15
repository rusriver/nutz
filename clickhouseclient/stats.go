package clickhouseclient

import (
	"encoding/json"

	"github.com/ClickHouse/clickhouse-go/v2/lib/proto"
	"github.com/dustin/go-humanize"
)

type Statistics struct {
	Status  string
	Elapsed string
	Read    struct {
		Rows      string
		Bytes     string
		RowsRate  string
		BytesRate string
	}
	Write struct {
		Rows      string
		Bytes     string
		RowsRate  string
		BytesRate string
	}
	progress proto.Progress
}

func (s *Statistics) Humanize() {
	s.Elapsed = (s.progress.Elapsed / 1e6 * 1e6).String()
	s.Read.Rows = humanize.SIWithDigits(float64(s.progress.Rows), 1, " rows")
	s.Read.Bytes = humanize.Bytes(s.progress.Bytes)
	s.Read.RowsRate = humanize.SIWithDigits(float64(s.progress.Rows)/s.progress.Elapsed.Seconds(), 1, " rows/s")
	s.Read.BytesRate = humanize.SIWithDigits(float64(s.progress.Bytes)/s.progress.Elapsed.Seconds(), 1, "B/s")
	s.Write.Rows = humanize.SIWithDigits(float64(s.progress.WroteRows), 1, " rows")
	s.Write.Bytes = humanize.Bytes(s.progress.WroteBytes)
	s.Write.RowsRate = humanize.SIWithDigits(float64(s.progress.WroteRows)/s.progress.Elapsed.Seconds(), 1, " rows/s")
	s.Write.BytesRate = humanize.SIWithDigits(float64(s.progress.WroteBytes)/s.progress.Elapsed.Seconds(), 1, "B/s")
}

func (s *Statistics) JSON() string {
	bb, _ := json.MarshalIndent(s, "", "    ")
	return string(bb)
}
