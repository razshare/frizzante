package receive

import (
	"github.com/razshare/frizzante/mock"
	"testing"
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
