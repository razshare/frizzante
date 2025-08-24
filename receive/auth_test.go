package receive

import (
	"encoding/base64"
	"github.com/razshare/frizzante/mock"
	"testing"
)

func TestBasicAuth(t *testing.T) {
	c := mock.NewClient()
	b := base64.URLEncoding.EncodeToString([]byte("test:123"))
	c.Request.Header.Set("Authorization", "Basic "+b)
	u, p, ok := BasicAuth(c)
	if !ok {
		t.Fatal("auth should pass")
	}
	if u != "test" {
		t.Fatal("user should be test")
	}
	if p != "123" {
		t.Fatal("password should be 123")
	}
}
