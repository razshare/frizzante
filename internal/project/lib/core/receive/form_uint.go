package receive

import (
	"strconv"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// FormUint reads the first form value associated with the given key
// as an uint and stores the result in the value pointed to by value.
//
// If there are no values associated with the key, FormUint stores 0 into value.
func FormUint(client *clients.Client, key string, value *uint) bool {
	var err error
	var valueLoc uint64
	text := FormValue(client, key)
	valueLoc, err = strconv.ParseUint(text, 10, 32)
	*value = uint(valueLoc)
	if err != nil {
		client.Config.ErrorLog.Println("form value is not a valid uint", stack.Trace())
		return false
	}
	return true
}

// FormUint32 reads the first form value associated with the given key
// as an uint32 and stores the result in the value pointed to by value.
//
// If there are no values associated with the key, FormUint32 stores 0 into value.
func FormUint32(client *clients.Client, key string, value *uint32) bool {
	var err error
	var valueLoc uint64
	text := FormValue(client, key)
	valueLoc, err = strconv.ParseUint(text, 10, 32)
	*value = uint32(valueLoc)
	if err != nil {
		client.Config.ErrorLog.Println("form value is not a valid uint32", stack.Trace())
		return false
	}
	return true
}

// FormUint64 reads the first form value associated with the given key
// as an uint64 and stores the result in the value pointed to by value.
//
// If there are no values associated with the key, FormUint64 stores 0 into value.
func FormUint64(client *clients.Client, key string, value *uint64) bool {
	var err error
	text := FormValue(client, key)
	*value, err = strconv.ParseUint(text, 10, 64)
	if err != nil {
		client.Config.ErrorLog.Println("form value is not a valid uint64", stack.Trace())
		return false
	}
	return true
}

// FormUintSlice reads all form values associated with the given key
// as a slice of uint and stores the result in the value pointed to by value.
//
// If there are no values associated with the key, FormUintSlice stores an empty slice into value.
// If any value fails to parse, FormUintSlice returns false.
func FormUintSlice(client *clients.Client, key string, value *[]uint) bool {
	texts := FormValues(client, key)
	*value = make([]uint, 0, len(texts))

	for _, text := range texts {
		var parsed uint64
		var err error
		if parsed, err = strconv.ParseUint(text, 10, 32); err != nil {
			client.Config.ErrorLog.Println("form value is not a valid uint", stack.Trace())
			return false
		}
		*value = append(*value, uint(parsed))
	}

	return true
}

// FormUint32Slice reads all form values associated with the given key
// as a slice of uint32 and stores the result in the value pointed to by value.
//
// If there are no values associated with the key, FormUint32Slice stores an empty slice into value.
// If any value fails to parse, FormUint32Slice returns false.
func FormUint32Slice(client *clients.Client, key string, value *[]uint32) bool {
	texts := FormValues(client, key)
	*value = make([]uint32, 0, len(texts))

	for _, text := range texts {
		var parsed uint64
		var err error
		if parsed, err = strconv.ParseUint(text, 10, 32); err != nil {
			client.Config.ErrorLog.Println("form value is not a valid uint32", stack.Trace())
			return false
		}
		*value = append(*value, uint32(parsed))
	}

	return true
}

// FormUint64Slice reads all form values associated with the given key
// as a slice of uint64 and stores the result in the value pointed to by value.
//
// If there are no values associated with the key, FormUint64Slice stores an empty slice into value.
// If any value fails to parse, FormUint64Slice returns false.
func FormUint64Slice(client *clients.Client, key string, value *[]uint64) bool {
	texts := FormValues(client, key)
	*value = make([]uint64, 0, len(texts))

	for _, text := range texts {
		var parsed uint64
		var err error
		if parsed, err = strconv.ParseUint(text, 10, 64); err != nil {
			client.Config.ErrorLog.Println("form value is not a valid uint64", stack.Trace())
			return false
		}
		*value = append(*value, parsed)
	}

	return true
}
