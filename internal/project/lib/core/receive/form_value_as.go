package receive

import (
	"strconv"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// FormValueAsFloat32 reads the first form value associated with the given key
// as a float32 and stores the result in the value pointed to by value.
//
// If there are no values associated with the key, FormValueAsFloat32 stores 0 into value.
func FormValueAsFloat32(client *clients.Client, key string, value *float32) bool {
	var err error
	var valueLoc float64
	text := FormValue(client, key)
	valueLoc, err = strconv.ParseFloat(text, 32)
	*value = float32(valueLoc)
	if err != nil {
		client.Config.ErrorLog.Println("form value is not a valid float32", stack.Trace())
		return false
	}
	return true
}

// FormValueAsFloat64 reads the first form value associated with the given key
// as a float64 and stores the result in the value pointed to by value.
//
// If there are no values associated with the key, FormValueAsFloat64 stores 0 into value.
func FormValueAsFloat64(client *clients.Client, key string, value *float64) bool {
	var err error
	text := FormValue(client, key)
	*value, err = strconv.ParseFloat(text, 64)
	if err != nil {
		client.Config.ErrorLog.Println("form value is not a valid float64", stack.Trace())
		return false
	}
	return true
}

// FormValueAsInt reads the first form value associated with the given key
// as an int and stores the result in the value pointed to by value.
//
// If there are no values associated with the key, FormValueAsInt stores 0 into value.
func FormValueAsInt(client *clients.Client, key string, value *int) bool {
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

// FormValueAsInt32 reads the first form value associated with the given key
// as an int32 and stores the result in the value pointed to by value.
//
// If there are no values associated with the key, FormValueAsInt32 stores 0 into value.
func FormValueAsInt32(client *clients.Client, key string, value *int32) bool {
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

// FormValueAsInt64 reads the first form value associated with the given key
// as an int64 and stores the result in the value pointed to by value.
//
// If there are no values associated with the key, FormValueAsInt64 stores 0 into value.
func FormValueAsInt64(client *clients.Client, key string, value *int64) bool {
	var err error
	text := FormValue(client, key)
	*value, err = strconv.ParseInt(text, 10, 64)
	if err != nil {
		client.Config.ErrorLog.Println("form value is not a valid int64", stack.Trace())
		return false
	}
	return true
}

// FormValueAsUint reads the first form value associated with the given key
// as an uint and stores the result in the value pointed to by value.
//
// If there are no values associated with the key, FormValueAsUint stores 0 into value.
func FormValueAsUint(client *clients.Client, key string, value *uint) bool {
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

// FormValueAsUint32 reads the first form value associated with the given key
// as an uint32 and stores the result in the value pointed to by value.
//
// If there are no values associated with the key, FormValueAsUint32 stores 0 into value.
func FormValueAsUint32(client *clients.Client, key string, value *uint32) bool {
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

// FormValueAsUint64 reads the first form value associated with the given key
// as an uint64 and stores the result in the value pointed to by value.
//
// If there are no values associated with the key, FormValueAsUint64 stores 0 into value.
func FormValueAsUint64(client *clients.Client, key string, value *uint64) bool {
	var err error
	text := FormValue(client, key)
	*value, err = strconv.ParseUint(text, 10, 64)
	if err != nil {
		client.Config.ErrorLog.Println("form value is not a valid uint64", stack.Trace())
		return false
	}
	return true
}
