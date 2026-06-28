package send

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestTooManyRequests(t *testing.T) {
	http := mocks.NewScope()
	TooManyRequests(http, "too many requests")
	writer := http.Writer.(*mocks.ResponseWriter)
	if http.Status != 429 {
		t.Fatal("status should be 429")
	}
	if string(writer.MockBytes) != "too many requests" {
		t.Fatal("content should be too many requests")
	}
}
