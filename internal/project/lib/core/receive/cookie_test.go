package receive

import (
	"main/lib/core/mock"
	"testing"
)

func TestCookie(t *testing.T) {
	client := mock.NewClient()
	client.Request.Header.Set("Cookie", "cookie=monster;")
	cookie := Cookie(client, "cookie")
	if cookie != "monster" {
		t.Fatal("cookie should be monster")
	}
}

func TestCookieEmptyKey(t *testing.T) {
	client := mock.NewClient()
	cookie := Cookie(client, "")
	if cookie != "" {
		t.Fatal("cookie should be empty")
	}
}

func TestCookieInvalidContent(t *testing.T) {
	client := mock.NewClient()
	client.Request.Header.Set("Cookie", "cookie=%monster;")
	cookie := Cookie(client, "cookie")
	if cookie != "" {
		t.Fatal("cookie should be empty")
	}
}
