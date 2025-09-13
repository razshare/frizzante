package send

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mock"
)

func TestJson(t *testing.T) {
	type Payload struct {
		Key string `json:"key"`
	}
	client := mock.NewClient()
	Json(client, Payload{Key: "value"})
	writer := client.Writer.(*mock.ResponseWriter)
	if string(writer.MockBytes) != `{"key":"value"}` {
		t.Fatal("content should be json")
	}
}
