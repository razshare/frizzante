package receive

import (
	"encoding/base64"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestBasicAuth(t *testing.T) {
	http := mocks.NewScope()
	token := base64.URLEncoding.EncodeToString([]byte("test:123"))
	http.Request.Header.Set("Authorization", "Basic "+token)
	username, password := BasicAuth(http)
	if username != "test" {
		t.Fatal("user should be test")
	}
	if password != "123" {
		t.Fatal("password should be 123")
	}
}
