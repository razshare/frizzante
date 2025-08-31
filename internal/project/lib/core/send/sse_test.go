package send

import (
	"main/lib/core/mock"
	"strings"
	"testing"
)

func TestSseUpgrade(t *testing.T) {
	client := mock.NewClient()
	SseUpgrade(client)
	if client.EventName != "message" {
		t.Fatal("event name should be message")
	}
}

func TestEventContentWithoutUpgrade(t *testing.T) {
	client := mock.NewClient()
	EventContent(client, []byte("hello"))
	writer := client.Writer.(*mock.ResponseWriter)
	if string(writer.MockBytes) != strings.Join([]string{"id: 1", "event: ", "data: hello", "", ""}, "\r\n") {
		t.Fatal("sse payload should contain data but not event name")
	}
}

func TestEventContentWithUpgrade(t *testing.T) {
	client := mock.NewClient()
	SseUpgrade(client)
	EventContent(client, []byte("hello"))
	writer := client.Writer.(*mock.ResponseWriter)
	if string(writer.MockBytes) != strings.Join([]string{"id: 1", "event: message", "data: hello", "", ""}, "\r\n") {
		t.Fatal("sse payload should contain event name message and data hello")
	}
}
