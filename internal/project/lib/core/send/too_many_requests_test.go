package send

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/dev/mocks"
)

func TestTooManyRequests(t *testing.T) {
	client := mocks.NewClient()
	TooManyRequests(client, "too many requests")
	writer := client.Writer.(*mocks.ResponseWriter)
	if client.Status != 429 {
		t.Fatal("status should be 429")
	}
	if string(writer.MockBytes) != "too many requests" {
		t.Fatal("content should be too many requests")
	}
}
