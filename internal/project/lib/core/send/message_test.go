package send

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/dev/mocks"
)

func TestMessage(t *testing.T) {
	client := mocks.NewClient()
	Message(client, "hello")
	writer := client.Writer.(*mocks.ResponseWriter)
	if string(writer.MockBytes) != "hello" {
		t.Fatal("content should be hello")
	}
}
