package send

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestMessagef(t *testing.T) {
	http := mocks.NewScope()
	Messagef(http, "hello %s", "world")
	writer := http.Writer.(*mocks.ResponseWriter)
	if string(writer.MockBytes) != "hello world" {
		t.Fatal("content should be hello world")
	}
}
