package genericmsboot

import (
	"github.com/rusriver/configdb/v3"
	loggeterrapi "github.com/rusriver/nutz/loggeterr-api"
)

type Glove struct {
	DiagnosticFlags *map[string]bool
	ConfDB          *configdb.Root
}

func (b *Boot) GetGloveAt(tracer loggeterrapi.ITracer, point_path string) *Glove {
	g := &Glove{
		DiagnosticFlags: b.DiagnosticFlags,
		ConfDB:          b.ConfDB.WithErr(nil).WithOk(nil).DoPanic().PP(point_path),
	}
	n := b.ConfDB.WithErr(nil).WithOk(nil).DoPanic()
	if tracer != nil {
		n = n.WithTracer(tracer)
	}
	g.ConfDB = n.PP(point_path)
	return g
}
