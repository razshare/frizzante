package send

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

// SseUpgrade upgrades to server sent events
// and returns a function that sets the name of the current event.
//
// The default event name is "message".
func SseUpgrade(client *clients.Client) func(event string) {
	Headers(client, map[string]string{
		"Access-Control-Allow-Origin":   "*",
		"Access-Control-Expose-Headers": "Content-Type",
		"Content-Type":                  "text/event-stream",
		"Cache-Control":                 "no-cache",
		"Client":                        "keep-alive",
	})

	client.EventName = "message"

	return func(event string) { client.EventName = event }
}

// EventContent sends content using the `server sent events` format.
//
// Usually this should be used internally in order to send content to a Server sent event.
//
// That being said, other than the format, there is nothing else different between this function and ResponseSendContent.
//
// See https://html.spec.whatwg.org/multipage/server-sent-events.html for more details on the format.
func EventContent(client *clients.Client, data []byte) {
	meta := fmt.Sprintf("id: %d\r\nevent: %s\r\n", client.EventId, client.EventName)
	if _, err := client.Writer.Write([]byte(meta)); err != nil {
		client.Options.ErrorLog.Println(err, stack.Trace())
		return
	}

	for _, line := range bytes.Split(data, []byte("\r\n")) {
		if _, err := client.Writer.Write([]byte("data: ")); err != nil {
			client.Options.ErrorLog.Println(err, stack.Trace())
			return
		}

		if _, err := client.Writer.Write(line); err != nil {
			client.Options.ErrorLog.Println(err, stack.Trace())
			return
		}

		if _, err := client.Writer.Write([]byte("\r\n")); err != nil {
			client.Options.ErrorLog.Println(err, stack.Trace())
			return
		}
	}

	if _, err := client.Writer.Write([]byte("\r\n")); err != nil {
		client.Options.ErrorLog.Println(err, stack.Trace())
		return
	}

	writer, ok := client.Writer.(http.Flusher)
	if !ok {
		client.Options.ErrorLog.Println(errors.New("could not retrieve flusher"), stack.Trace())
		return
	}

	writer.Flush()

	client.EventId++
}
