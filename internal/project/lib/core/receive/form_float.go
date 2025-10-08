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
//
// If value fails to parse, FormInts returns false.
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
