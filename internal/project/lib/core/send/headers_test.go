package send

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestHeaders(t *testing.T) {
	client := mocks.NewClient()
	Headers(client, map[string]string{
		"key1": "value1",
		"key2": "value2",
	})
	writer := client.Writer.(*mocks.ResponseWriter)
	if writer.MockHeader.Get("key1") != "value1" {
		t.Fatal("key1 should be value1")
	}
	if writer.MockHeader.Get("key2") != "value2" {
		t.Fatal("key2 should be value2")
	}
}
