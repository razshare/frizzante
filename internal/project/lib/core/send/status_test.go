package send

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/dev/mocks"
)

func TestStatus(t *testing.T) {
	client := mocks.NewClient()
	Status(client, 400)
	if client.Status != 400 {
		t.Fatal("status should be 400")
	}
}
