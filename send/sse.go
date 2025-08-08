package send

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/razshare/frizzante/conn"
	"github.com/razshare/frizzante/stack"
	"net/http"
)

// EventContent sends content using the `server sent events` format.
//
// Usually this should be used internally in order to send content to a Server sent event.
//
// That being said, other than the format, there is nothing else different between this function and ResponseSendContent.
//
// See https://html.spec.whatwg.org/multipage/server-sent-events.html for more details on the format.
func EventContent(c *conn.Conn, d []byte) {
	header := fmt.Sprintf("id: %d\r\nevent: %s\r\n", c.EventId, c.EventName)

	_, writeError := c.Writer.Write([]byte(header))
	if writeError != nil {
		c.Container.Config.ErrorLog.Println(writeError, stack.Trace())
		return
	}

	for _, line := range bytes.Split(d, []byte("\r\n")) {
		_, writeError = c.Writer.Write([]byte("data: "))
		if writeError != nil {
			c.Container.Config.ErrorLog.Println(writeError, stack.Trace())
			return
		}

		_, writeError = c.Writer.Write(line)
		if writeError != nil {
			c.Container.Config.ErrorLog.Println(writeError, stack.Trace())
			return
		}

		_, writeError = c.Writer.Write([]byte("\r\n"))
		if writeError != nil {
			c.Container.Config.ErrorLog.Println(writeError, stack.Trace())
			return
		}
	}

	_, writeError = c.Writer.Write([]byte("\r\n"))
	if writeError != nil {
		c.Container.Config.ErrorLog.Println(writeError, stack.Trace())
		return
	}

	flusher, flushedOk := c.Writer.(http.Flusher)
	if !flushedOk {
		c.Container.Config.ErrorLog.Println(errors.New("could not retrieve flusher"), stack.Trace())
		return
	}

	flusher.Flush()

	c.EventId++
}

// SseUpgrade upgrades to server sent events
// and returns a function that sets the name of the current event.
//
// The default event name is "message".
func SseUpgrade(c *conn.Conn) func(string) {
	Headers(c, map[string]string{
		"Access-Control-Allow-Origin":   "*",
		"Access-Control-Expose-Headers": "Content-Type",
		"Content-Type":                  "text/event-stream",
		"Cache-Control":                 "no-cache",
		"Conn":                          "keep-alive",
	})

	c.EventName = "message"

	return func(eventName string) { c.EventName = eventName }
}
