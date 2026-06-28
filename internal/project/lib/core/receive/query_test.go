package receive

import (
	"net/url"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestQuery(t *testing.T) {
	http := mocks.NewScope()
	http.Request.URL = &url.URL{RawQuery: "key1=value1&key2=value2"}
	if Query(http, "key1") != "value1" {
		t.Fatal("key1 should be value1")
	}
	if Query(http, "key2") != "value2" {
		t.Fatal("key2 should be value2")
	}
}
