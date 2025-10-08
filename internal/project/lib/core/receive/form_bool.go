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
func FormBool(client *clients.Client, key string) bool {
	var err error
	var value bool
	text := FormValue(client, key)
	value, err = strconv.ParseBool(text)
	if err != nil {
		client.Config.ErrorLog.Println("form value is not a valid bool", stack.Trace())
		return false
	}
	return value
}
