package send

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestBadRequest(t *testing.T) {
	http := mocks.NewScope()
	BadRequest(http, "bad request")
	writer := http.Writer.(*mocks.ResponseWriter)
	if http.Status != 400 {
		t.Fatal("status should be 400")
	}
	if string(writer.MockBytes) != "bad request" {
		t.Fatal("content should be bad request")
	}
}
