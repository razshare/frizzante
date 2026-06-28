package negotiate

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestSessionId(t *testing.T) {
	http := mocks.NewScope()
	http.Request.Header.Set("Cookie", "session-id=value;")
	if SessionId(http) != "value" {
		t.Fatal("session id should be value")
	}
}

func TestSessionIdCached(t *testing.T) {
	http := mocks.NewScope()
	http.SessionId = "value"
	if SessionId(http) != "value" {
		t.Fatal("session id should be value")
	}
}
