package dslre

type ICache[T any] interface {
	Get(key string) (val T, ok bool)
	Set(key string, val T)
}
