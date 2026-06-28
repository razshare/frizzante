package send

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestTooManyRequestsf(t *testing.T) {
	http := mocks.NewScope()
	TooManyRequestsf(http, "too many %s", "requests")
	writer := http.Writer.(*mocks.ResponseWriter)
	if http.Status != 429 {
		t.Fatal("status should be 429")
	}
	if string(writer.MockBytes) != "too many requests" {
		t.Fatal("content should be too many requests")
	}
}
