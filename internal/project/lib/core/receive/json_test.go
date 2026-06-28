package receive

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestJson(t *testing.T) {
	type Payload struct {
		Key string `json:"key"`
	}
	http := mocks.NewScope()
	body := http.Request.Body.(*mocks.RequestBody)
	body.MockBuffer = []byte(`{"key":"value"}`)
	var payload Payload
	Json(http, &payload)
	if payload.Key != "value" {
		t.Fatal("key should be value")
	}
}
