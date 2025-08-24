package receive

import (
	"github.com/razshare/frizzante/mock"
	"net/url"
	"testing"
)

func TestQuery(t *testing.T) {
	c := mock.NewClient()
	c.Request.URL = &url.URL{RawQuery: "key1=value1&key2=value2"}
	if Query(c, "key1") != "value1" {
		t.Fatal("key1 should be value1")
	}
	if Query(c, "key2") != "value2" {
		t.Fatal("key2 should be value2")
	}
}
