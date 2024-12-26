package controlflow

type KeyedTable[T comparable] map[T]func() int // -1 break/complete

func (t KeyedTable[T]) Sequence(keys ...T) {
	for _, casKey := range keys {
		if then, ok := t[casKey]; ok {
			switch then() {
			case -1:
				return
			}
		}
	}
}

func (t KeyedTable[T]) Switch(primary, default_ T) {
	if then, ok := t[primary]; ok {
		then()
	} else {
		if then, ok := t[default_]; ok {
			then()
		}
	}
}
