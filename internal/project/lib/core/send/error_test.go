package send

import (
	"errors"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestError(t *testing.T) {
	http := mocks.NewScope()
	Error(http, errors.New("error"))
	writer := http.Writer.(*mocks.ResponseWriter)
	if http.Status != 500 {
		t.Fatal("status should be 500")
	}
	if string(writer.MockBytes) != "error" {
		t.Fatal("content should be error")
	}
}
