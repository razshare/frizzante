package receive

import (
	"errors"
	"net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/client"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// Parse reads the message as a form.
//
// The whole request body is parsed and up to a total of MaxFormSize bytes
// of its file parts are stored in memory, with the remainder stored on disk in temporary files.
func Parse(client *client.Client) {
	if client.Parsed {
		return
	}

	client.Parsed = true

	if client.WebSocket != nil {
		client.Config.ErrorLog.Println("web socket connections cannot parse forms", stack.Trace())
		return
	}

	if err := client.Request.ParseMultipartForm(MaxFormSize); err != nil {
		if !errors.Is(err, http.ErrNotMultipart) {
			return
		}

		err = client.Request.ParseForm()
		if err != nil {
			client.Config.ErrorLog.Println(err, stack.Trace())
			return
		}
	}
}
