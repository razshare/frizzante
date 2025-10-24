package receive

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestPath(t *testing.T) {
	client := mocks.NewClient()
	client.Request.SetPathValue("key", "value")
	if Path(client, "key") != "value" {
		t.Fatal("key should be value")
	}
}
