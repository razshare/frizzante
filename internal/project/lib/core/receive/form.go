package receive

import (
	"mime/multipart"
	"reflect"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
)

func Form(client *clients.Client, value any) bool {
	reflection := reflect.ValueOf(value)

	if reflection.Kind() == reflect.Pointer {
		reflection = reflection.Elem()
	}

	type_ := reflection.Type()
	for index := range reflection.NumField() {
		var ok bool
		var name string
		reflectionValue := reflection.Field(index)
		reflectionField := type_.Field(index)

		if tag := reflectionField.Tag.Get("form"); tag != "" {
			name = tag
		} else {
			name = reflectionField.Name
		}

		switch reflectionValue.Interface().(type) {
		case string:
			var local string
			ok = FormValue(client, name, &local)
			reflectionValue.Set(reflect.ValueOf(local))
		case []byte:
			var local []byte
			ok = FormValue(client, name, &local)
			reflectionValue.Set(reflect.ValueOf(local))
		case bool:
			var local bool
			ok = FormValue(client, name, &local)
			reflectionValue.Set(reflect.ValueOf(local))
		case []bool:
			var local []bool
			ok = FormValue(client, name, &local)
			reflectionValue.Set(reflect.ValueOf(local))
		case uint:
			var local uint
			ok = FormValue(client, name, &local)
			reflectionValue.Set(reflect.ValueOf(local))
		case []uint:
			var local []uint
			ok = FormValue(client, name, &local)
			reflectionValue.Set(reflect.ValueOf(local))
		case uint32:
			var local uint32
			ok = FormValue(client, name, &local)
			reflectionValue.Set(reflect.ValueOf(local))
		case []uint32:
			var local []uint32
			ok = FormValue(client, name, &local)
			reflectionValue.Set(reflect.ValueOf(local))
		case uint64:
			var local uint64
			ok = FormValue(client, name, &local)
			reflectionValue.Set(reflect.ValueOf(local))
		case []uint64:
			var local []uint64
			ok = FormValue(client, name, &local)
			reflectionValue.Set(reflect.ValueOf(local))
		case int:
			var local int
			ok = FormValue(client, name, &local)
			reflectionValue.Set(reflect.ValueOf(local))
		case []int:
			var local []int
			ok = FormValue(client, name, &local)
			reflectionValue.Set(reflect.ValueOf(local))
		case int32:
			var local int32
			ok = FormValue(client, name, &local)
			reflectionValue.Set(reflect.ValueOf(local))
		case []int32:
			var local []int32
			ok = FormValue(client, name, &local)
			reflectionValue.Set(reflect.ValueOf(local))
		case int64:
			var local int64
			ok = FormValue(client, name, &local)
			reflectionValue.Set(reflect.ValueOf(local))
		case []int64:
			var local []int64
			ok = FormValue(client, name, &local)
			reflectionValue.Set(reflect.ValueOf(local))
		case float32:
			var local float32
			ok = FormValue(client, name, &local)
			reflectionValue.Set(reflect.ValueOf(local))
		case []float32:
			var local []float32
			ok = FormValue(client, name, &local)
			reflectionValue.Set(reflect.ValueOf(local))
		case float64:
			var local float64
			ok = FormValue(client, name, &local)
			reflectionValue.Set(reflect.ValueOf(local))
		case []float64:
			var local []float64
			ok = FormValue(client, name, &local)
			reflectionValue.Set(reflect.ValueOf(local))
		case multipart.FileHeader:
			var local multipart.FileHeader
			ok = FormValue(client, name, &local)
			reflectionValue.Set(reflect.ValueOf(local))
		case []multipart.FileHeader:
			var local []multipart.FileHeader
			ok = FormValue(client, name, &local)
			reflectionValue.Set(reflect.ValueOf(local))
		case []*multipart.FileHeader:
			var local []*multipart.FileHeader
			ok = FormValue(client, name, &local)
			reflectionValue.Set(reflect.ValueOf(local))
		case multipart.File:
			var local multipart.File
			ok = FormValue(client, name, &local)
			reflectionValue.Set(reflect.ValueOf(local))
		case []multipart.File:
			var local []multipart.File
			ok = FormValue(client, name, &local)
			reflectionValue.Set(reflect.ValueOf(local))
		case []*multipart.File:
			var local []*multipart.File
			ok = FormValue(client, name, &local)
			reflectionValue.Set(reflect.ValueOf(local))
		}
		if !ok {
			return false
		}
	}

	return true
}
