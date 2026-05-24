package send

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestForbiddenf(t *testing.T) {
	client := mocks.NewClient()
	Forbiddenf(client, "%s", "forbidden")
	writer := client.Writer.(*mocks.ResponseWriter)
	if client.Status != 403 {
		t.Fatal("status should be 403")
	}
	if string(writer.MockBytes) != "forbidden" {
		t.Fatal("content should be forbidden")
	}
}
