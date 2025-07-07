package genericmsboot

import "github.com/rusriver/configdb/v3"

type Glove struct {
	DiagnosticFlags map[string]bool
	ConfDB          *configdb.Root
}

func (b *Boot) GetGloveAt(point_path string) *Glove {
	g := &Glove{
		DiagnosticFlags: b.DiagnosticFlags,
		ConfDB:          b.ConfDB.DoPanic().PP(point_path),
	}
	return g
}
