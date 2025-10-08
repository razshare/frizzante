package receive

import (
	"strconv"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// FormInts reads all form values associated with the given key
// as a slice of int and stores the result in the value pointed to by value.
//
// If there are no values associated with the key, FormInts stores an empty slice into value.
//
// If any value fails to parse, FormInts returns false.
func FormInts(client *clients.Client, key string, value *[]int) bool {
	texts := FormValues(client, key)
	*value = make([]int, len(texts))

	for index, text := range texts {
		var parsed int64
		var err error
		if parsed, err = strconv.ParseInt(text, 10, 32); err != nil {
			client.Config.ErrorLog.Println("form value is not a valid int", stack.Trace())
			return false
		}
		(*value)[index] = int(parsed)
	}

	return true
}

// FormInt32s reads all form values associated with the given key
// as a slice of int32 and stores the result in the value pointed to by value.
//
// If there are no values associated with the key, FormInt32s stores an empty slice into value.
//
// If any value fails to parse, FormInt32s returns false.
func FormInt32s(client *clients.Client, key string, value *[]int32) bool {
	texts := FormValues(client, key)
	*value = make([]int32, 0, len(texts))

	for index, text := range texts {
		var parsed int64
		var err error
		if parsed, err = strconv.ParseInt(text, 10, 32); err != nil {
			client.Config.ErrorLog.Println("form value is not a valid int32", stack.Trace())
			return false
		}
		(*value)[index] = int32(parsed)
	}

	return true
}

// FormInt64s reads all form values associated with the given key
// as a slice of int64 and stores the result in the value pointed to by value.
//
// If there are no values associated with the key, FormInt64s stores an empty slice into value.
//
// If any value fails to parse, FormInt64s returns false.
func FormInt64s(client *clients.Client, key string, value *[]int64) bool {
	texts := FormValues(client, key)
	*value = make([]int64, 0, len(texts))

	for index, text := range texts {
		var parsed int64
		var err error
		if parsed, err = strconv.ParseInt(text, 10, 64); err != nil {
			client.Config.ErrorLog.Println("form value is not a valid int64", stack.Trace())
			return false
		}
		(*value)[index] = parsed
	}

	return true
}
