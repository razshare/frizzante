package types

import "reflect"

// EncodeValue encodes a value recursively.
//
// Slices and arrays are encoded as slices. Complex types within the resulting slices are also encoded.
//
// Interfaces, structs and maps are encoded as maps. Complex types within the resulting maps are also encoded.
//
// Pointers are unwrapped.
func EncodeValue(value reflect.Value) (out any, err error) {
	kind := value.Kind()
	switch kind {
	case reflect.Slice,
		reflect.Array:
		count := value.Len()
		items := make([]any, count)
		for i := 0; i < count; i++ {
			item := value.Index(i)
			if items[i], err = EncodeValue(item); err != nil {
				return
			}
		}
		out = items
	case reflect.Struct,
		reflect.Map,
		reflect.Interface:
		out, err = EncodeInterface(value.Interface())
	case reflect.Pointer:
		out, err = EncodeValue(value.Elem())
	default:
		out = value.Interface()
	}
	return
}
