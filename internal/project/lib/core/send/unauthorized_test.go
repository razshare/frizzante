package send

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/dev/mocks"
)

func TestUnauthorized(t *testing.T) {
	client := mocks.NewClient()
	Unauthorized(client, "unauthorized")
	writer := client.Writer.(*mocks.ResponseWriter)

	if client.Status != 401 {
		t.Fatal("status should be 401")
	}

	if string(writer.MockBytes) != "unauthorized" {
		t.Fatal("content should be unauthorized")
	}
}
