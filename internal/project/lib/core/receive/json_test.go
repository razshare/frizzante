package receive

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestJson(t *testing.T) {
	type Payload struct {
		Key string `json:"key"`
	}
	request, _ := mocks.NewExchange()
	body := request.Body.(*mocks.RequestBody)
	body.MockBuffer = []byte(`{"key":"value"}`)
	var payload Payload
	if err := Json(request, &payload); err != nil {
		t.Fatal("json parsing should succeed")
	}
	if payload.Key != "value" {
		t.Fatal("key should be value")
	}
}
