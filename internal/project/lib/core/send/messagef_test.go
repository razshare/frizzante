package send

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestMessagef(t *testing.T) {
	client := mocks.NewClient()
	Messagef(client, "hello %s", "world")
	writer := client.Writer.(*mocks.ResponseWriter)

	if string(writer.MockBytes) != "hello world" {
		t.Fatal("content should be hello world")
	}
}
