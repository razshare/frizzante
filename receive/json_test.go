package receive

import (
	"github.com/razshare/frizzante/mock"
	"testing"
)

func TestJson(t *testing.T) {
	type Payload struct {
		Key string `json:"key"`
	}
	c := mock.NewClient()
	b := c.Request.Body.(*mock.RequestBody)
	b.MockBuffer = []byte(`{"key":"value"}`)
	p := Json[Payload](c)
	if p.Key != "value" {
		t.Fatal("key should be value")
	}
}
