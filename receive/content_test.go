package receive

import (
	"github.com/razshare/frizzante/mock"
	"io"
	"testing"
)

func TestMessage(t *testing.T) {
	c := mock.NewClient()
	b := c.Request.Body.(*mock.RequestBody)
	b.MockBuffer = []byte("hello")
	d, err := io.ReadAll(c.Request.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(d) != "hello" {
		t.Fatal("request body should be hello")
	}
}
