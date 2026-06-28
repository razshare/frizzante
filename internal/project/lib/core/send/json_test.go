package send

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestJson(t *testing.T) {
	type Payload struct {
		Key string `json:"key"`
	}
	_, writer := mocks.NewExchange()
	if err := Json(writer, Payload{Key: "value"}); err != nil {
		t.Fatal("json encoding should succeed")
	}
	if string(writer.MockBytes) != `{"key":"value"}` {
		t.Fatal("content should be json")
	}
}
