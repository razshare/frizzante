package receive

import (
	"strconv"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// FormInt reads the first form value associated with the given key
// as an int and stores the result in the value pointed to by value.
//
// If there are no values associated with the key, FormInt stores 0 into value.
func FormInt(client *clients.Client, key string, value *int) bool {
	var err error
	var valueLoc int64
	text := FormValue(client, key)
	valueLoc, err = strconv.ParseInt(text, 10, 32)
	*value = int(valueLoc)
	if err != nil {
		client.Config.ErrorLog.Println("form value is not a valid int", stack.Trace())
		return false
	}
	return true
}

// FormInt32 reads the first form value associated with the given key
// as an int32 and stores the result in the value pointed to by value.
//
// If there are no values associated with the key, FormInt32 stores 0 into value.
func FormInt32(client *clients.Client, key string, value *int32) bool {
	var err error
	var valueLoc int64
	text := FormValue(client, key)
	valueLoc, err = strconv.ParseInt(text, 10, 32)
	*value = int32(valueLoc)
	if err != nil {
		client.Config.ErrorLog.Println("form value is not a valid int32", stack.Trace())
		return false
	}
	return true
}

// FormInt64 reads the first form value associated with the given key
// as an int64 and stores the result in the value pointed to by value.
//
// If there are no values associated with the key, FormInt64 stores 0 into value.
func FormInt64(client *clients.Client, key string, value *int64) bool {
	var err error
	text := FormValue(client, key)
	*value, err = strconv.ParseInt(text, 10, 64)
	if err != nil {
		client.Config.ErrorLog.Println("form value is not a valid int64", stack.Trace())
		return false
	}
	return true
}

// FormIntSlice reads all form values associated with the given key
// as a slice of int and stores the result in the value pointed to by value.
//
// If there are no values associated with the key, FormIntSlice stores an empty slice into value.
// If any value fails to parse, FormIntSlice returns false.
func FormIntSlice(client *clients.Client, key string, value *[]int) bool {
	texts := FormValues(client, key)
	*value = make([]int, 0, len(texts))

	for _, text := range texts {
		var parsed int64
		var err error
		if parsed, err = strconv.ParseInt(text, 10, 32); err != nil {
			client.Config.ErrorLog.Println("form value is not a valid int", stack.Trace())
			return false
		}
		*value = append(*value, int(parsed))
	}

	return true
}

// FormInt32Slice reads all form values associated with the given key
// as a slice of int32 and stores the result in the value pointed to by value.
//
// If there are no values associated with the key, FormInt32Slice stores an empty slice into value.
// If any value fails to parse, FormInt32Slice returns false.
func FormInt32Slice(client *clients.Client, key string, value *[]int32) bool {
	texts := FormValues(client, key)
	*value = make([]int32, 0, len(texts))

	for _, text := range texts {
		var parsed int64
		var err error
		if parsed, err = strconv.ParseInt(text, 10, 32); err != nil {
			client.Config.ErrorLog.Println("form value is not a valid int32", stack.Trace())
			return false
		}
		*value = append(*value, int32(parsed))
	}

	return true
}

// FormInt64Slice reads all form values associated with the given key
// as a slice of int64 and stores the result in the value pointed to by value.
//
// If there are no values associated with the key, FormInt64Slice stores an empty slice into value.
// If any value fails to parse, FormInt64Slice returns false.
func FormInt64Slice(client *clients.Client, key string, value *[]int64) bool {
	texts := FormValues(client, key)
	*value = make([]int64, 0, len(texts))

	for _, text := range texts {
		var parsed int64
		var err error
		if parsed, err = strconv.ParseInt(text, 10, 64); err != nil {
			client.Config.ErrorLog.Println("form value is not a valid int64", stack.Trace())
			return false
		}
		*value = append(*value, parsed)
	}

	return true
}
