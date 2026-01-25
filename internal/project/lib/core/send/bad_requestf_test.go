package send

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/dev/mocks"
)

func TestBadRequestf(t *testing.T) {
	client := mocks.NewClient()
	BadRequestf(client, "bad %s", "request")
	writer := client.Writer.(*mocks.ResponseWriter)
	if client.Status != 400 {
		t.Fatal("status should be 400")
	}
	if string(writer.MockBytes) != "bad request" {
		t.Fatal("content should be bad request")
	}
}
