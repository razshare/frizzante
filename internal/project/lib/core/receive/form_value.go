package receive

import (
	"errors"
	"net/http"
	"strconv"

	_client "github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// FormValue reads the first form value associated with the given key and returns it.
//
// If there are no values associated with the key, FormValue returns an empty string.
func FormValue(client *_client.Client, key string) string {
	if client.WebSocket != nil {
		client.Config.ErrorLog.Println("web socket connections cannot parse forms", stack.Trace())
		return ""
	}

	if client.Request.Form == nil {
		if err := client.Request.ParseMultipartForm(MaxFormSize); err != nil {
			if !errors.Is(err, http.ErrNotMultipart) {
				return ""
			}
		}
	}

	return client.Request.Form.Get(key)
}

// FormValueIsFloat checks if the form value associated with the given key is a valid float.
func FormValueIsFloat(client *_client.Client, key string) bool {
	value := FormValue(client, key)
	if value == "" {
		return false
	}
	_, err := strconv.ParseFloat(value, 64)
	return err == nil
}

// FormValueAsFloat reads the first form value associated with the given key and returns it as a float64.
//
// If there are no values associated with the key or the value is not a valid float, FormValueAsFloat returns 0.
//
// Use FormValueIsFloat to make sure the key exists and its value is a float.
func FormValueAsFloat(client *_client.Client, key string) float64 {
	value := FormValue(client, key)
	if value == "" {
		return 0
	}

	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		client.Config.ErrorLog.Println("form value is not a valid float", stack.Trace())
		return 0
	}

	return result
}
