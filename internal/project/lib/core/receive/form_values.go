package receive

import (
	"errors"
	"net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// FormValues reads all form values associated with the given key and returns them as a slice.
//
// If there are no values associated with the key, FormValues returns an empty slice.
func FormValues(client *clients.Client, key string) []string {
	if client.WebSocket != nil {
		client.Config.ErrorLog.Println("web socket connections cannot parse forms", stack.Trace())
		return make([]string, 0)
	}

	if client.Request.Form == nil {
		if err := client.Request.ParseMultipartForm(MaxFormSize); err != nil {
			if !errors.Is(err, http.ErrNotMultipart) {
				return make([]string, 0)
			}
		}
	}

	if values, ok := client.Request.Form[key]; ok {
		return values
	}

	return make([]string, 0)
}
