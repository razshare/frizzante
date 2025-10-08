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
