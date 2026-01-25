package send

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/dev/mocks"
)

func TestContent(t *testing.T) {
	client := mocks.NewClient()
	Content(client, []byte("hello"))
	writer := client.Writer.(*mocks.ResponseWriter)
	if string(writer.MockBytes) != "hello" {
		t.Fatal("content should be hello")
	}
}
