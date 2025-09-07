package simplecache

type Cache[T any] struct {
	Map map[string]T
}

func NewCache[T any]() (c *Cache[T]) {
	c = &Cache[T]{
		Map: make(map[string]T),
	}
	return c
}

func (c *Cache[T]) Get(key string) (val T, ok bool) {
	val, ok = c.Map[key]
	return
}

func (c *Cache[T]) Set(key string, val T) {
	c.Map[key] = val
}
