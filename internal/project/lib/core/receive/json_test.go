package receive

import (
	"main/lib/core/mock"
	"testing"
)

func TestJson(t *testing.T) {
	type Payload struct {
		Key string `json:"key"`
	}
	client := mock.NewClient()
	body := client.Request.Body.(*mock.RequestBody)
	body.MockBuffer = []byte(`{"key":"value"}`)
	var payload Payload
	Json(client, &payload)
	if payload.Key != "value" {
		t.Fatal("key should be value")
	}
}
