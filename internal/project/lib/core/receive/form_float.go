package receive

import (
	"strconv"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// FormFloat32 reads the first form value associated with the given key
// as a float32 and stores the result in the value pointed to by value.
//
// If there are no values associated with the key, FormFloat32 stores 0 into value.
func FormFloat32(client *clients.Client, key string, value *float32) bool {
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

// FormFloat64 reads the first form value associated with the given key
// as a float64 and stores the result in the value pointed to by value.
//
// If there are no values associated with the key, FormFloat64 stores 0 into value.
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
