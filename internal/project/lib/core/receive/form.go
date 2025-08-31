package receive

import (
	"errors"
	"main/lib/core/client"
	"net/http"
	"net/url"

	"github.com/razshare/frizzante/globals"
	"github.com/razshare/frizzante/stack"
)

// Form reads the message as a form and returns the value.
//
// The whole request body is parsed and up to a total of 2MB
// of its file parts are stored in memory, with the remainder stored on disk in temporary files.
func Form(client *client.Client) url.Values {
	return FormWithMaxMemory(client, 2*globals.MB)
}

// FormWithMaxMemory reads the message as a form and returns the value.
//
// The whole request body is parsed and up to a total of m bytes
// of its file parts are stored in memory, with the remainder stored on disk in temporary files.
func FormWithMaxMemory(client *client.Client, max int64) url.Values {
	if client.WebSocket != nil {
		client.Config.ErrorLog.Println("web socket connections cannot parse forms", stack.Trace())
		return url.Values{}
	}

	if err := client.Request.ParseMultipartForm(max); err != nil {
		if !errors.Is(err, http.ErrNotMultipart) {
			return url.Values{}
		}

		err = client.Request.ParseForm()
		if err != nil {
			client.Config.ErrorLog.Println(err, stack.Trace())
			return url.Values{}
		}
	}

	return client.Request.Form
}
