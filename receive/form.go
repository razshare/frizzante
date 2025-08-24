package receive

import (
	"errors"
	"github.com/razshare/frizzante/client"
	"github.com/razshare/frizzante/globals"
	"github.com/razshare/frizzante/stack"
	"net/http"
	"net/url"
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
func FormWithMaxMemory(c *client.Client, m int64) url.Values {
	if c.WebSocket != nil {
		c.Config.ErrorLog.Println(errors.New("web socket connections cannot parse forms"), stack.Trace())
		return url.Values{}
	}

	if err := c.Request.ParseMultipartForm(m); err != nil {
		if !errors.Is(err, http.ErrNotMultipart) {
			return url.Values{}
		}

		err = c.Request.ParseForm()
		if err != nil {
			c.Config.ErrorLog.Println(err, stack.Trace())
			return url.Values{}
		}
	}

	return c.Request.Form
}
