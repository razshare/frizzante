package receive

import (
	"strconv"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// FormUints reads all form values associated with the given key
// as a slice of uint and stores the result in the value pointed to by value.
//
// If there are no values associated with the key, FormUints stores an empty slice into value.
//
// If any value fails to parse, FormUints returns false.
func FormUints(client *clients.Client, key string, value *[]uint) bool {
	texts := FormValues(client, key)
	*value = make([]uint, len(texts))

	for index, text := range texts {
		var parsed uint64
		var err error
		if parsed, err = strconv.ParseUint(text, 10, 32); err != nil {
			client.Config.ErrorLog.Println("form value is not a valid uint", stack.Trace())
			return false
		}
		(*value)[index] = uint(parsed)
	}

	return true
}

// FormUint32s reads all form values associated with the given key
// as a slice of uint32 and stores the result in the value pointed to by value.
//
// If there are no values associated with the key, FormUint32s stores an empty slice into value.
//
// If any value fails to parse, FormUint32s returns false.
func FormUint32s(client *clients.Client, key string, value *[]uint32) bool {
	texts := FormValues(client, key)
	*value = make([]uint32, len(texts))

	for index, text := range texts {
		var parsed uint64
		var err error
		if parsed, err = strconv.ParseUint(text, 10, 32); err != nil {
			client.Config.ErrorLog.Println("form value is not a valid uint32", stack.Trace())
			return false
		}
		(*value)[index] = uint32(parsed)
	}

	return true
}

// FormUint64s reads all form values associated with the given key
// as a slice of uint64 and stores the result in the value pointed to by value.
//
// If there are no values associated with the key, FormUint64s stores an empty slice into value.
//
// If any value fails to parse, FormUint64s returns false.
func FormUint64s(client *clients.Client, key string, value *[]uint64) bool {
	texts := FormValues(client, key)
	*value = make([]uint64, len(texts))

	for index, text := range texts {
		var parsed uint64
		var err error
		if parsed, err = strconv.ParseUint(text, 10, 64); err != nil {
			client.Config.ErrorLog.Println("form value is not a valid uint64", stack.Trace())
			return false
		}
		(*value)[index] = parsed
	}

	return true
}
