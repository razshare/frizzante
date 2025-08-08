package receive

import (
	"encoding/json"
	"errors"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/razshare/frizzante/conn"
	"github.com/razshare/frizzante/globals"
	"github.com/razshare/frizzante/send"
	"github.com/razshare/frizzante/stack"
	"io"
	"net/http"
	"net/url"
)

// SessionId tries to find a session id among the user's cookies.
// If no session id is found, it creates a new one and returns it.
func SessionId(c *conn.Conn) string {
	var id string
	cookies := c.Request.CookiesNamed("session-id")
	cookiesCount := 0

	for _, cookie := range cookies {
		id = cookie.Value
		cookiesCount++
	}

	if cookiesCount > 0 {
		return id
	}

	// Create new session.
	idObject, idObjectError := uuid.NewV4()
	if idObjectError != nil {
		c.Container.Config.ErrorLog.Println(idObjectError, stack.Trace())
		return ""
	}

	id = idObject.String()

	send.Cookie(c, "session-id", id)

	return id
}

// Cancellation returns a channel that closes when the request gets cancelled.
func Cancellation(c *conn.Conn) <-chan struct{} {
	return c.Request.Context().Done()
}

// IsAlive returns a reference to a bool which is initially set to `true`.
//
// This bool updates to `false` when the request gets cancelled.
func IsAlive(c *conn.Conn) *bool {
	isAlive := true
	go func() {
		<-Cancellation(c)
		isAlive = false
	}()
	return &isAlive
}

// Cookie reads the contents of a cookie from the message and returns the value.
//
// Compatible with web sockets.
func Cookie(c *conn.Conn, key string) string {
	cookie, cookieError := c.Request.Cookie(key)
	if cookieError != nil {
		c.Container.Config.ErrorLog.Println(cookieError, stack.Trace())
		return ""
	}

	data, queryError := url.QueryUnescape(cookie.Value)
	if queryError != nil {
		c.Container.Config.ErrorLog.Println(queryError, stack.Trace())
		return ""
	}

	return data
}

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

// Query reads a query field and returns the value.
//
// Compatible with web sockets.
func Query(c *conn.Conn, k string) string {
	return c.Request.URL.Query().Get(k)
}

// Path reads a parameters fields and returns the value.
//
// Compatible with web sockets.
func Path(c *conn.Conn, k string) string {
	return c.Request.PathValue(k)
}

// Header reads a header field and returns the value.
//
// Compatible with web sockets.
func Header(c *conn.Conn, k string) string {
	return c.Request.Header.Get(k)
}

// ContentType reads the Content-Type header field and returns the value.
//
// Compatible with web sockets.
func ContentType(c *conn.Conn) string {
	return c.Request.Header.Get("Content-Type")
}

// Accept reads if the Accept header entries and returns the values.
func Accept(c *conn.Conn) string {
	return c.Request.Header.Get("Accept")
}
