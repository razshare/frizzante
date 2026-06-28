package send

import (
	"github.com/razshare/frizzante/internal/project/lib/core/scopes"
)

// SseUpgrade upgrades to server sent events
// and returns a function that sets the name of the current event.
//
// The default event name is "message".
func SseUpgrade(http *scopes.Http) func(event string) {
	Headers(http, map[string]string{
		"Access-Control-Expose-Headers": "Content-Type",
		"Content-Type":                  "text/event-stream",
		"Cache-Control":                 "no-cache",
		"Client":                        "keep-alive",
	})
	http.EventName = "message"
	return func(event string) { http.EventName = event }
}
