package receive

import (
	"strconv"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// FormBool reads the first form value associated with the given key
// as a bool and stores the result in the value pointed to by value.
//
// Valid bool values are: 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False.
//
// If there are no values associated with the key, FormBool stores false into value.
//
// If value fails to parse, FormInts returns false.
func FormBool(client *clients.Client, key string, value *bool) bool {
	var err error
	text := FormValue(client, key)
	if *value, err = strconv.ParseBool(text); err != nil {
		client.Config.ErrorLog.Println("form value is not a valid bool", stack.Trace())
		return false
	}
	return true
}
