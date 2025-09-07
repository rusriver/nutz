package examples_test

import (
	"embed"

	"github.com/rusriver/nutz/dslre"
)

//go:embed test-literals-storage/*
var Storage0 embed.FS

var Storage dslre.IReadFile

type TStorage struct {
	Prefix string
	embed.FS
}

func (s *TStorage) ReadFile(name string) ([]byte, error) {
	return s.FS.ReadFile(s.Prefix + name)
}

func init() {
	Storage = &TStorage{Prefix: "test-literals-storage/", FS: Storage0}
}
