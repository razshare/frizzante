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
func Form(c *client.Client) url.Values {
	return FormWithMaxMemory(c, 2*globals.MB)
}

// FormWithMaxMemory reads the message as a form and returns the value.
//
// The whole request body is parsed and up to a total of maxMemory bytes
// of its file parts are stored in memory, with the remainder stored on disk in temporary files.
func FormWithMaxMemory(c *client.Client, m int64) url.Values {
	if c.WebSocket != nil {
		c.Config.ErrorLog.Println(errors.New("web socket connections cannot parse forms"), stack.Trace())
		return url.Values{}
	}

	formError := c.Request.ParseMultipartForm(m)
	if formError != nil {
		if !errors.Is(formError, http.ErrNotMultipart) {
			return url.Values{}
		}

		formError = c.Request.ParseForm()
		if formError != nil {
			c.Config.ErrorLog.Println(formError, stack.Trace())
			return url.Values{}
		}
	}

	return c.Request.Form
}
