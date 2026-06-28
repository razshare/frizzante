package send

import (
	"strings"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestEventContentWithoutUpgrade(t *testing.T) {
	http := mocks.NewScope()
	EventContent(http, []byte("hello"))
	writer := http.Writer.(*mocks.ResponseWriter)
	if string(writer.MockBytes) != strings.Join([]string{"id: 1", "event: ", "data: hello", "", ""}, "\r\n") {
		t.Fatal("sse payload should contain data but not event name")
	}
}

func TestEventContentWithUpgrade(t *testing.T) {
	http := mocks.NewScope()
	SseUpgrade(http)
	EventContent(http, []byte("hello"))
	writer := http.Writer.(*mocks.ResponseWriter)
	if string(writer.MockBytes) != strings.Join([]string{"id: 1", "event: message", "data: hello", "", ""}, "\r\n") {
		t.Fatal("sse payload should contain event name message and data hello")
	}
}
