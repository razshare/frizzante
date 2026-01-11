package send

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestNotFound(t *testing.T) {
	client := mocks.NewClient()
	NotFound(client, "not found")
	writer := client.Writer.(*mocks.ResponseWriter)

	if client.Status != 404 {
		t.Fatal("status should be 404")
	}

	if string(writer.MockBytes) != "not found" {
		t.Fatal("content should be not found")
	}
}
