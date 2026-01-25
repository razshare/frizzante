package send

import (
	"github.com/razshare/frizzante/internal/project/lib/core/clients"
)

// SseUpgrade upgrades to server sent events
// and returns a function that sets the name of the current event.
//
// The default event name is "message".
func SseUpgrade(client *clients.Client) func(event string) {
	Headers(client, map[string]string{
		"Access-Control-Expose-Headers": "Content-Type",
		"Content-Type":                  "text/event-stream",
		"Cache-Control":                 "no-cache",
		"Client":                        "keep-alive",
	})
	client.EventName = "message"
	return func(event string) { client.EventName = event }
}
