package send

import (
	"github.com/razshare/frizzante/mock"
	"testing"
)

func TestJson(t *testing.T) {
	type Payload struct {
		Key string `json:"key"`
	}
	c := mock.NewClient()
	Json(c, Payload{Key: "value"})
	w := c.Writer.(*mock.ResponseWriter)
	if string(w.MockBytes) != `{"key":"value"}` {
		t.Fatal("content should be json")
	}
}
