package receive

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestSessionId(t *testing.T) {
	client := mocks.NewClient()
	client.Request.Header.Set("Cookie", "session-id=value;")
	if SessionId(client) != "value" {
		t.Fatal("session id should be value")
	}
}

func TestSessionIdCached(t *testing.T) {
	client := mocks.NewClient()
	client.SessionId = "value"
	if SessionId(client) != "value" {
		t.Fatal("session id should be value")
	}
}
