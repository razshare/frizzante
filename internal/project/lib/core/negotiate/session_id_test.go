package negotiate

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestSessionId(t *testing.T) {
	request, writer := mocks.NewExchange()
	request.Header.Set("Cookie", "session-id=value;")
	if value, err := SessionId(writer, request); err != nil || value != "value" {
		t.Fatal("session id should be value")
	}
}
