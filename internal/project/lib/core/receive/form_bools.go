package receive

import (
	"strconv"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// FormBools reads all form values associated with the given key
// as a slice of bool and stores the result in the value pointed to by value.
//
// Valid bool values are: 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False.
//
// If there are no values associated with the key, FormBools stores an empty slice into value.
//
// If any value fails to parse, FormBools returns false.
func FormBools(client *clients.Client, key string, value *[]bool) bool {
	texts := FormValues(client, key)
	*value = make([]bool, len(texts))

	for index, text := range texts {
		var parsed bool
		var err error
		if parsed, err = strconv.ParseBool(text); err != nil {
			client.Config.ErrorLog.Println("form value is not a valid bool", stack.Trace())
			return false
		}
		(*value)[index] = parsed
	}

	return true
}
