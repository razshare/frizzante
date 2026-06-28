package send

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestSseUpgrade(t *testing.T) {
	http := mocks.NewScope()
	SseUpgrade(http)
	if http.EventName != "message" {
		t.Fatal("event name should be message")
	}
}
