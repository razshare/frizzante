package types

import (
	"errors"
	"reflect"
	"strings"
)

type Article struct {
	Title string
}

// EncodeValue encodes a value recursively.
//
// Slices and arrays are encoded as slices. Complex types within the resulting slices are also encoded.
//
// Interfaces, structs and maps are encoded as maps. Complex types within the resulting maps are also encoded.
//
// Pointers are unwrapper.
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

// EncodeInterface encodes a map or a struct as a deep map.
func EncodeInterface(from any) (out map[string]any, err error) {
	if from == nil {
		out = nil
		return
	}

	value := reflect.ValueOf(from)
	kind := value.Kind()

	if kind == reflect.Pointer {
		value = value.Elem()
	}

	if kind == reflect.Map {
		out = from.(map[string]any)
		return
	}

	if kind != reflect.Struct {
		return nil, errors.New("only maps and structs can be converted to deep maps")
	}

	type_ := value.Type()
	count := value.NumField()
	out = make(map[string]any)
	for i := 0; i < count; i++ {
		field := type_.Field(i)
		if strings.ToLower(field.Name[0:1]) == field.Name[0:1] {
			continue
		}

		if json := field.Tag.Get("json"); json != "" {
			if out[json], err = EncodeValue(value.Field(i)); err != nil {
				return
			}
		} else {
			if out[field.Name], err = EncodeValue(value.Field(i)); err != nil {
				return
			}
		}
	}
	return out, nil
}
