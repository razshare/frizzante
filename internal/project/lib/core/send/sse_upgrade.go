package send

import (
	"net/http"
)

// SseUpgrade upgrades to server sent events
// and returns a function that sets the name of the current event.
//
// The default event name is "message".
func SseUpgrade(writer *http.ResponseWriter) func(event string) {
	headers := (*writer).Header()
	headers.Set("Access-Control-Expose-Headers", "Content-Type")
	headers.Set("Content-Type", "text/event-stream")
	headers.Set("Cache-Control", "no-cache")
	headers.Set("Client", "keep-alive")
	upgrade := &SseUpgradeResponseWriter{
		ResponseWriter: *writer,
		EventName:      "message",
		EventId:        1,
	}
	converted := http.ResponseWriter(upgrade)
	*writer = converted
	return func(event string) {
		upgrade.EventName = event
	}
}
