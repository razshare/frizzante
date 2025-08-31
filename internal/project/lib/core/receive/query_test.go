package receive

import (
	"main/lib/core/mock"
	"net/url"
	"testing"
)

func TestQuery(t *testing.T) {
	client := mock.NewClient()
	client.Request.URL = &url.URL{RawQuery: "key1=value1&key2=value2"}
	if Query(client, "key1") != "value1" {
		t.Fatal("key1 should be value1")
	}
	if Query(client, "key2") != "value2" {
		t.Fatal("key2 should be value2")
	}
}
