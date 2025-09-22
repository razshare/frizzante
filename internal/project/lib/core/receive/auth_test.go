package receive

import (
	"encoding/base64"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mock"
)

func TestBasicAuth(t *testing.T) {
	client := mock.NewClient()
	token := base64.URLEncoding.EncodeToString([]byte("test:123"))
	client.Request.Header.Set("Authorization", "Basic "+token)

	username, password := BasicAuth(client)

	if username != "test" {
		t.Fatal("user should be test")
	}
	if password != "123" {
		t.Fatal("password should be 123")
	}
}
