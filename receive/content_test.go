package receive

import (
	"github.com/razshare/frizzante/mock"
	"io"
	"testing"
)

func TestMessage(t *testing.T) {
	client := mock.NewClient()
	body := client.Request.Body.(*mock.RequestBody)
	body.MockBuffer = []byte("hello")
	data, err := io.ReadAll(client.Request.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "hello" {
		t.Fatal("request body should be hello")
	}
}
