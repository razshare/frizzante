package send

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestMessage(t *testing.T) {
	http := mocks.NewScope()
	Message(http, "hello")
	writer := http.Writer.(*mocks.ResponseWriter)
	if string(writer.MockBytes) != "hello" {
		t.Fatal("content should be hello")
	}
}
