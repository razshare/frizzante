package send

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/razshare/frizzante/client"
	"github.com/razshare/frizzante/stack"
	"net/http"
)

// SseUpgrade upgrades to server sent events
// and returns a function that sets the name of the current event.
//
// The default event name is "message".
func SseUpgrade(c *client.Client) func(string) {
	Headers(c, map[string]string{
		"Access-Control-Allow-Origin":   "*",
		"Access-Control-Expose-Headers": "Content-Type",
		"Content-Type":                  "text/event-stream",
		"Cache-Control":                 "no-cache",
		"Client":                        "keep-alive",
	})

	c.EventName = "message"

	return func(n string) { c.EventName = n }
}

// EventContent sends content using the `server sent events` format.
//
// Usually this should be used internally in order to send content to a Server sent event.
//
// That being said, other than the format, there is nothing else different between this function and ResponseSendContent.
//
// See https://html.spec.whatwg.org/multipage/server-sent-events.html for more details on the format.
func EventContent(c *client.Client, d []byte) {
	h := fmt.Sprintf("id: %d\r\nevent: %s\r\n", c.EventId, c.EventName)

	if _, err := c.Writer.Write([]byte(h)); err != nil {
		c.Config.ErrorLog.Println(err, stack.Trace())
		return
	}

	for _, l := range bytes.Split(d, []byte("\r\n")) {
		if _, err := c.Writer.Write([]byte("data: ")); err != nil {
			c.Config.ErrorLog.Println(err, stack.Trace())
			return
		}

		if _, err := c.Writer.Write(l); err != nil {
			c.Config.ErrorLog.Println(err, stack.Trace())
			return
		}

		if _, err := c.Writer.Write([]byte("\r\n")); err != nil {
			c.Config.ErrorLog.Println(err, stack.Trace())
			return
		}
	}

	if _, err := c.Writer.Write([]byte("\r\n")); err != nil {
		c.Config.ErrorLog.Println(err, stack.Trace())
		return
	}

	w, ok := c.Writer.(http.Flusher)
	if !ok {
		c.Config.ErrorLog.Println(errors.New("could not retrieve flusher"), stack.Trace())
		return
	}

	w.Flush()

	c.EventId++
}
