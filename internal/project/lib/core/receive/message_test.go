package receive

import (
	"io"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestMessage(t *testing.T) {
	http := mocks.NewScope()
	body := http.Request.Body.(*mocks.RequestBody)
	body.MockBuffer = []byte("hello")
	data, err := io.ReadAll(http.Request.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "hello" {
		t.Fatal("request body should be hello")
	}
}
