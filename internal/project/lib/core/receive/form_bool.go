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
func FormBool(client *clients.Client, key string, value *bool) bool {
	var err error
	text := FormValue(client, key)
	if *value, err = strconv.ParseBool(text); err != nil {
		client.Config.ErrorLog.Println("form value is not a valid bool", stack.Trace())
		return false
	}
	return true
}

// FormBoolSlice reads all form values associated with the given key
// as a slice of bool and stores the result in the value pointed to by value.
//
// Valid bool values are: 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False.
//
// If there are no values associated with the key, FormBoolSlice stores an empty slice into value.
// If any value fails to parse, FormBoolSlice returns false.
func FormBoolSlice(client *clients.Client, key string, value *[]bool) bool {
	texts := FormValues(client, key)
	*value = make([]bool, 0, len(texts))

	for _, text := range texts {
		var parsed bool
		var err error
		if parsed, err = strconv.ParseBool(text); err != nil {
			client.Config.ErrorLog.Println("form value is not a valid bool", stack.Trace())
			return false
		}
		*value = append(*value, parsed)
	}

	return true
}
