package receive

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/dev/mocks"
)

func TestJson(t *testing.T) {
	type Payload struct {
		Key string `json:"key"`
	}
	client := mocks.NewClient()
	body := client.Request.Body.(*mocks.RequestBody)
	body.MockBuffer = []byte(`{"key":"value"}`)
	var payload Payload
	Json(client, &payload)
	if payload.Key != "value" {
		t.Fatal("key should be value")
	}
}
