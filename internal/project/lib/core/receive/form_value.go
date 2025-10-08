package receive

import (
	"errors"
	"net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// FormValue reads the first form value associated with the given key and returns it.
//
// If there are no values associated with the key, FormValue returns an empty string.
func FormValue(client *clients.Client, key string) string {
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

// FormValues reads all form values associated with the given key and returns them as a slice.
//
// If there are no values associated with the key, FormValues returns an empty slice.
//
// This is the array equivalent of FormValue, which returns only the first value.
func FormValues(client *clients.Client, key string) []string {
	if client.WebSocket != nil {
		client.Config.ErrorLog.Println("web socket connections cannot parse forms", stack.Trace())
		return []string{}
	}

	if client.Request.Form == nil {
		if err := client.Request.ParseMultipartForm(MaxFormSize); err != nil {
			if !errors.Is(err, http.ErrNotMultipart) {
				return []string{}
			}
		}
	}

	if values, ok := client.Request.Form[key]; ok {
		return values
	}

	return []string{}
}
