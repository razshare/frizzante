package receive

import (
	"github.com/razshare/frizzante/mock"
	"testing"
)

func TestSessionId(t *testing.T) {
	c := mock.NewClient()
	c.Request.Header.Set("Cookie", "session-id=value;")
	if SessionId(c) != "value" {
		t.Fatal("session id should be value")
	}
}

func TestSessionIdCached(t *testing.T) {
	c := mock.NewClient()
	c.SessionId = "value"
	if SessionId(c) != "value" {
		t.Fatal("session id should be value")
	}
}
