package receive

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mock"
)

func TestSessionId(t *testing.T) {
	client := mock.NewClient()
	client.Request.Header.Set("Cookie", "session-id=value;")
	if SessionId(client) != "value" {
		t.Fatal("session id should be value")
	}
}

func TestSessionIdCached(t *testing.T) {
	client := mock.NewClient()
	client.SessionId = "value"
	if SessionId(client) != "value" {
		t.Fatal("session id should be value")
	}
}
