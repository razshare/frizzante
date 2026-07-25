package receive

import (
	"testing"

	"github.com/razshare/frizzante/v2/internal/project/lib/core/mocks"
)

func TestCookie(t *testing.T) {
	request, _ := mocks.NewExchange()
	request.Header.Set("Cookie", "cookie=monster;")

	if cookie, err := Cookie(request, "cookie"); err != nil || cookie != "monster" {
		t.Fatal("cookie should be monster")
	}
}

func TestCookieEmptyKey(t *testing.T) {
	request, _ := mocks.NewExchange()
	if cookie, err := Cookie(request, ""); err == nil || cookie != "" {
		t.Fatal("cookie should be empty")
	}
}

func TestCookieInvalidContent(t *testing.T) {
	request, _ := mocks.NewExchange()
	request.Header.Set("Cookie", "cookie=%monster;")
	if cookie, err := Cookie(request, "cookie"); err == nil || cookie != "" {
		t.Fatal("cookie should be empty")
	}
}
