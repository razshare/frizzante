package send

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestJson(t *testing.T) {
	type Payload struct {
		Key string `json:"key"`
	}
	http := mocks.NewScope()
	Json(http, Payload{Key: "value"})
	writer := http.Writer.(*mocks.ResponseWriter)
	if string(writer.MockBytes) != `{"key":"value"}` {
		t.Fatal("content should be json")
	}
}
