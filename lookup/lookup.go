/*
Small library on top of reflect for make lookups to Structs or Maps. Using a
very simple DSL you can access to any property, key or value of any value of Go.
*/
package lookup

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrInvalidIndexUsage = errors.New("invalid index usage")
	ErrKeyNotFound       = errors.New("not found")
	ErrIndexOutOfBounds  = errors.New("index out of bounds")
)

// DotP performs a lookup into a value, using a string. Same as `Lookup`
// but using a string with the keys separated by `.`
func DotP(i interface{}, path string) (reflect.Value, error) {
	return P(i, strings.Split(path, ".")...)
}

// P performs a lookup into a value, using a path of keys. The key should
// match with a Field or a MapIndex. For slice you can use the syntax key[index]
// to access a specific index. If one key owns to a slice and an index is not
// specificied the rest of the path will be apllied to evaley value of the
// slice, and the value will be merged into a slice.
func P(i interface{}, path ...string) (reflect.Value, error) {
	return lookup(i, path...)
}

func lookup(i interface{}, path ...string) (reflect.Value, error) {
	nodeValue := reflect.ValueOf(i)
	var err error
	for _, part := range path {
		nodeValue, err = getChildByName(nodeValue, part)
		if err == nil {
			continue
		}
		break
	}
	return nodeValue, err
}

func getChildByName(nodeValue reflect.Value, key string) (childValue reflect.Value, err error) {

	nodeKind := nodeValue.Kind()
	// fmt.Println("++53", nodeKind, key)
	switch nodeKind {

	case reflect.Ptr, reflect.Interface:
		return getChildByName(nodeValue.Elem(), key)

	case reflect.Struct:
		childValue = nodeValue.FieldByName(key)

	case reflect.Map:
		kValue := reflect.Indirect(reflect.New(nodeValue.Type().Key()))
		kValue.SetString(key)
		childValue = nodeValue.MapIndex(kValue)

	case reflect.Array, reflect.Slice:
		index, er := strconv.Atoi(key)
		if er != nil {
			err = ErrInvalidIndexUsage
			return
		}
		if nodeValue.Len() <= index {
			err = ErrIndexOutOfBounds
			return
		}
		childValue = nodeValue.Index(index)

	}
	if !childValue.IsValid() {
		err = ErrKeyNotFound
		return
	}
	if childValue.Kind() == reflect.Ptr || childValue.Kind() == reflect.Interface {
		childValue = childValue.Elem()
	}
	return childValue, nil
}
