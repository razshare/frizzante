package receive

import (
	"io"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestMessage(t *testing.T) {
	client := mocks.NewClient()
	body := client.Request.Body.(*mocks.RequestBody)
	body.MockBuffer = []byte("hello")
	data, err := io.ReadAll(client.Request.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "hello" {
		t.Fatal("request body should be hello")
	}
}
