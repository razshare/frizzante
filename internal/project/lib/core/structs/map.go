package structs

import (
	"reflect"
)

func Map(from any) map[string]any {
	result := make(map[string]any)

	reflection := reflect.ValueOf(from)

	if reflection.Kind() == reflect.Ptr {
		reflection = reflection.Elem()
	}

	reflectionType := reflection.Type()

	for i := 0; i < reflection.NumField(); i++ {
		structField := reflectionType.Field(i)

		var key string
		if jsonName := structField.Tag.Get("json"); jsonName != "" {
			key = jsonName
		} else {
			key = structField.Name
		}

		value := reflection.Field(i)

		if value.Kind() == reflect.Struct {
			result[key] = Map(value.Interface())
		} else {
			result[key] = value.Interface()
		}
	}

	return result
}
