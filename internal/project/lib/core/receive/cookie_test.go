package receive

import (
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestCookie(t *testing.T) {
	http := mocks.NewScope()
	http.Request.Header.Set("Cookie", "cookie=monster;")
	cookie := Cookie(http, "cookie")
	if cookie != "monster" {
		t.Fatal("cookie should be monster")
	}
}

func TestCookieEmptyKey(t *testing.T) {
	http := mocks.NewScope()
	cookie := Cookie(http, "")
	if cookie != "" {
		t.Fatal("cookie should be empty")
	}
}

func TestCookieInvalidContent(t *testing.T) {
	http := mocks.NewScope()
	http.Request.Header.Set("Cookie", "cookie=%monster;")
	cookie := Cookie(http, "cookie")
	if cookie != "" {
		t.Fatal("cookie should be empty")
	}
}
