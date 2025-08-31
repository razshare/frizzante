package receive

import (
	"main/lib/core/mock"
	"testing"
)

func TestPath(t *testing.T) {
	client := mock.NewClient()
	client.Request.SetPathValue("key", "value")
	if Path(client, "key") != "value" {
		t.Fatal("key should be value")
	}
}
