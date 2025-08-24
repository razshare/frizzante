package send

import (
	"github.com/razshare/frizzante/mock"
	"strings"
	"testing"
)

func TestSseUpgrade(t *testing.T) {
	c := mock.NewClient()
	SseUpgrade(c)
	if c.EventName != "message" {
		t.Fatal("event name should be message")
	}
}

func TestEventContentWithoutUpgrade(t *testing.T) {
	c := mock.NewClient()
	EventContent(c, []byte("hello"))
	w := c.Writer.(*mock.ResponseWriter)
	if string(w.MockBytes) != strings.Join([]string{"id: 1", "event: ", "data: hello", "", ""}, "\r\n") {
		t.Fatal("sse payload should contain data but not event name")
	}
}

func TestEventContentWithUpgrade(t *testing.T) {
	c := mock.NewClient()
	SseUpgrade(c)
	EventContent(c, []byte("hello"))
	w := c.Writer.(*mock.ResponseWriter)
	if string(w.MockBytes) != strings.Join([]string{"id: 1", "event: message", "data: hello", "", ""}, "\r\n") {
		t.Fatal("sse payload should contain event name message and data hello")
	}
}
