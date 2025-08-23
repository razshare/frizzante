package send

import (
	"strings"
	"testing"
)

func TestSseUpgrade(t *testing.T) {
	client := MockClient()
	SseUpgrade(client)
	if client.EventName != "message" {
		t.Fatal("event name should be message")
	}
}

func TestEventContentWithoutUpgrade(t *testing.T) {
	client := MockClient()
	EventContent(client, []byte("hello"))
	writer := client.Writer.(*MockWriter)
	actual := string(writer.MockBytes)
	expected := strings.Join(
		[]string{
			"id: 1",
			"event: ",
			"data: hello",
			"",
			"",
		},
		"\r\n",
	)
	if actual != expected {
		t.Fatal("sse ")
	}
}

func TestEventContentWithUpgrade(t *testing.T) {
	client := MockClient()
	SseUpgrade(client)
	EventContent(client, []byte("hello"))
	writer := client.Writer.(*MockWriter)
	actual := string(writer.MockBytes)
	expected := strings.Join(
		[]string{
			"id: 1",
			"event: message",
			"data: hello",
			"",
			"",
		},
		"\r\n",
	)
	if actual != expected {
		t.Fatal("sse ")
	}
}
