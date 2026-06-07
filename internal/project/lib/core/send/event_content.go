package send

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/razshare/frizzante/internal/project/lib/core/clients"
	"github.com/razshare/frizzante/internal/project/lib/core/logs"
	"github.com/razshare/frizzante/internal/project/lib/core/stack"
)

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
		logs.Errorf(
			client,
			"send.EventContent: failed to write event meta: %v\n%s",
			err,
			stack.Trace(),
		)
		return
	}
	for _, line := range bytes.Split(data, []byte("\r\n")) {
		if _, err := client.Writer.Write([]byte("data: ")); err != nil {
			logs.Errorf(
				client,
				"send.EventContent: failed to write data prefix: %v\n%s",
				err,
				stack.Trace(),
			)
			return
		}
		if _, err := client.Writer.Write(line); err != nil {
			logs.Errorf(
				client,
				"send.EventContent: failed to write data line: %v\n%s",
				err,
				stack.Trace(),
			)
			return
		}
		if _, err := client.Writer.Write([]byte("\r\n")); err != nil {
			logs.Errorf(
				client,
				"send.EventContent: failed to write line ending: %v\n%s",
				err,
				stack.Trace(),
			)
			return
		}
	}
	if _, err := client.Writer.Write([]byte("\r\n")); err != nil {
		logs.Errorf(
			client,
			"send.EventContent: failed to write event ending: %v\n%s",
			err,
			stack.Trace(),
		)
		return
	}
	writer, ok := client.Writer.(http.Flusher)
	if !ok {
		logs.Errorf(
			client,
			"send.EventContent: could not retrieve flusher\n%s",
			stack.Trace(),
		)
		return
	}
	writer.Flush()
	client.EventId++
}
