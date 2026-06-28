package receive

import (
	"io"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestMessage(t *testing.T) {
	request, _ := mocks.NewExchange()
	body := request.Body.(*mocks.RequestBody)
	body.MockBuffer = []byte("hello")
	data, err := io.ReadAll(request.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello" {
		t.Fatal("request body should be hello")
	}
}
