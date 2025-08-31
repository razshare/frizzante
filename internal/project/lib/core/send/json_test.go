package send

import (
	"main/lib/core/mock"
	"testing"
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
