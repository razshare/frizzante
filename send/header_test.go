package send

import (
	"github.com/razshare/frizzante/mock"
	"testing"
)

func TestHeader(t *testing.T) {
	c := mock.NewClient()
	Header(c, "key", "value")
	w := c.Writer.(*mock.ResponseWriter)
	if w.MockHeader.Get("key") != "value" {
		t.Fatal("key should be value")
	}
}

func TestHeaders(t *testing.T) {
	c := mock.NewClient()
	Headers(c, map[string]string{
		"key1": "value1",
		"key2": "value2",
	})

	w := c.Writer.(*mock.ResponseWriter)

	if w.MockHeader.Get("key1") != "value1" {
		t.Fatal("key1 should be value1")
	}

	if w.MockHeader.Get("key2") != "value2" {
		t.Fatal("key2 should be value2")
	}
}

func TestRedirect(t *testing.T) {
	c := mock.NewClient()
	Redirect(c, "/about", 303)
	w := c.Writer.(*mock.ResponseWriter)

	if c.Status != 303 {
		t.Fatal("status should be 303")
	}

	l := w.MockHeader.Get("Location")

	if l != "/about" {
		t.Fatal("location should be about")
	}
}

func TestNavigate(t *testing.T) {
	c := mock.NewClient()
	Navigate(c, "/about")
	w := c.Writer.(*mock.ResponseWriter)

	if c.Status != 302 {
		t.Fatal("status should be 302")
	}

	l := w.MockHeader.Get("Location")

	if l != "/about" {
		t.Fatal("location should be about")
	}
}

func TestContentType(t *testing.T) {
	c := mock.NewClient()
	ContentType(c, "text/html")
	w := c.Writer.(*mock.ResponseWriter)

	ct := w.MockHeader.Get("Content-Type")

	if ct != "text/html" {
		t.Fatal("content type should be text/html")
	}
}
