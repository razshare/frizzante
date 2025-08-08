package receive

import (
	"encoding/json"
	"errors"
	"github.com/razshare/frizzante/conn"
	"github.com/razshare/frizzante/globals"
	"github.com/razshare/frizzante/stack"
	"io"
	"net/http"
	"net/url"
)

// Message reads the contents of the message and returns the value.
//
// Compatible with web sockets.
func Message(c *conn.Conn) string {
	if c.WebSocket != nil {
		_, data, readError := c.WebSocket.ReadMessage()
		if readError != nil {
			c.Container.Config.ErrorLog.Println(readError, stack.Trace())
			return ""
		}
		return string(data)
	}

	data, readError := io.ReadAll(c.Request.Body)
	if readError != nil {
		c.Container.Config.ErrorLog.Println(readError, stack.Trace())
		return ""
	}
	return string(data)
}

// Json reads the next JSON-encoded message from the
// c and stores it in the value pointed to by va.
//
// Compatible with web sockets.
func Json(c *conn.Conn, v any) {
	if c.WebSocket != nil {
		jsonError := c.WebSocket.ReadJSON(v)
		if jsonError != nil {
			c.Container.Config.ErrorLog.Println(jsonError, stack.Trace())
			return
		}
		return
	}

	data, readError := io.ReadAll(c.Request.Body)
	if readError != nil {
		c.Container.Config.ErrorLog.Println(readError, stack.Trace())
		return
	}

	jsonError := json.Unmarshal(data, v)
	if jsonError != nil {
		c.Container.Config.ErrorLog.Println(jsonError, stack.Trace())
		return
	}
}

// Form reads the message as a form and returns the value.
//
// The whole request body is parsed and up to a total of 2MB
// of its file parts are stored in memory, with the remainder stored on disk in temporary files.
func Form(c *conn.Conn) url.Values {
	return FormWithMaxMemory(c, 2*globals.MB)
}

// FormWithMaxMemory reads the message as a form and returns the value.
//
// The whole request body is parsed and up to a total of maxMemory bytes
// of its file parts are stored in memory, with the remainder stored on disk in temporary files.
func FormWithMaxMemory(c *conn.Conn, m int64) url.Values {
	if c.WebSocket != nil {
		c.Container.Config.ErrorLog.Println(errors.New("c is not of type web socket"), stack.Trace())
		return url.Values{}
	}

	formError := c.Request.ParseMultipartForm(m)
	if formError != nil {
		if !errors.Is(formError, http.ErrNotMultipart) {
			return url.Values{}
		}

		formError = c.Request.ParseForm()
		if formError != nil {
			c.Container.Config.ErrorLog.Println(formError, stack.Trace())
			return url.Values{}
		}
	}

	return c.Request.Form
}
