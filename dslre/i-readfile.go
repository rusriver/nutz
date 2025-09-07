package dslre

type IReadFile interface {
	ReadFile(name string) ([]byte, error)
}
