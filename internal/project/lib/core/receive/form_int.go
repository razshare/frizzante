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
