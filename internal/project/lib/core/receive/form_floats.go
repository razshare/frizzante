package receive

import (
	"strconv"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// FormFloat64 reads the first form value associated with the given key
// as a float64 and stores the result in the value pointed to by value.
//
// If there are no values associated with the key, FormFloat64 stores 0 into value.
//
// If value fails to parse, FormInts returns false.
func FormFloat64(client *clients.Client, key string, value *float64) bool {
	var err error
	text := FormValue(client, key)
	*value, err = strconv.ParseFloat(text, 64)
	if err != nil {
		client.Config.ErrorLog.Println("form value is not a valid float64", stack.Trace())
		return false
	}
	return true
}

// FormFloat32s reads all form values associated with the given key
// as a slice of float32 and stores the result in the value pointed to by value.
//
// If there are no values associated with the key, FormFloat32s stores an empty slice into value.
//
// If any value fails to parse, FormFloat32s returns false.
func FormFloat32s(client *clients.Client, key string, value *[]float32) bool {
	texts := FormValues(client, key)
	*value = make([]float32, len(texts))

	for index, text := range texts {
		var parsed float64
		var err error
		if parsed, err = strconv.ParseFloat(text, 32); err != nil {
			client.Config.ErrorLog.Println("form value is not a valid float32", stack.Trace())
			return false
		}
		(*value)[index] = float32(parsed)
	}

	return true
}

// FormFloat64s reads all form values associated with the given key
// as a slice of float64 and stores the result in the value pointed to by value.
//
// If there are no values associated with the key, FormFloat64s stores an empty slice into value.
//
// If any value fails to parse, FormFloat64s returns false.
func FormFloat64s(client *clients.Client, key string, value *[]float64) bool {
	texts := FormValues(client, key)
	*value = make([]float64, len(texts))

	for index, text := range texts {
		var parsed float64
		var err error
		if parsed, err = strconv.ParseFloat(text, 64); err != nil {
			client.Config.ErrorLog.Println("form value is not a valid float64", stack.Trace())
			return false
		}
		(*value)[index] = parsed
	}

	return true
}
