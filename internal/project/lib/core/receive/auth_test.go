package receive

import (
	"encoding/base64"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mock"
)

func TestBasicAuth(t *testing.T) {
	client := mock.NewClient()
	text64 := base64.URLEncoding.EncodeToString([]byte("test:123"))
	client.Request.Header.Set("Authorization", "Basic "+text64)
	user, pass, ok := BasicAuth(client)
	if !ok {
		t.Fatal("auth should pass")
	}
	if user != "test" {
		t.Fatal("user should be test")
	}
	if pass != "123" {
		t.Fatal("password should be 123")
	}
}
