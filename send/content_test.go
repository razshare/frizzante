package send

import (
	"errors"
	"github.com/razshare/frizzante/mock"
	"testing"
)

func TestContent(t *testing.T) {
	c := mock.NewClient()
	Content(c, []byte("hello"))
	w := c.Writer.(*mock.ResponseWriter)

	if string(w.MockBytes) != "hello" {
		t.Fatal("content should be hello")
	}
}

func TestMessage(t *testing.T) {
	c := mock.NewClient()
	Message(c, "hello")
	w := c.Writer.(*mock.ResponseWriter)

	if string(w.MockBytes) != "hello" {
		t.Fatal("content should be hello")
	}
}

func TestMessagef(t *testing.T) {
	c := mock.NewClient()
	Messagef(c, "hello %s", "world")
	w := c.Writer.(*mock.ResponseWriter)

	if string(w.MockBytes) != "hello world" {
		t.Fatal("content should be hello world")
	}
}

func TestNotFound(t *testing.T) {
	c := mock.NewClient()
	NotFound(c, "not found")
	w := c.Writer.(*mock.ResponseWriter)

	if c.Status != 404 {
		t.Fatal("status should be 404")
	}

	if string(w.MockBytes) != "not found" {
		t.Fatal("content should be not found")
	}
}

func TestUnauthorized(t *testing.T) {
	c := mock.NewClient()
	Unauthorized(c, "unauthorized")
	w := c.Writer.(*mock.ResponseWriter)

	if c.Status != 401 {
		t.Fatal("status should be 401")
	}

	if string(w.MockBytes) != "unauthorized" {
		t.Fatal("content should be unauthorized")
	}
}

func TestBadRequest(t *testing.T) {
	c := mock.NewClient()
	BadRequest(c, "bad request")
	w := c.Writer.(*mock.ResponseWriter)

	if c.Status != 400 {
		t.Fatal("status should be 400")
	}

	if string(w.MockBytes) != "bad request" {
		t.Fatal("content should be bad request")
	}
}

func TestError(t *testing.T) {
	c := mock.NewClient()
	Error(c, errors.New("error"))
	w := c.Writer.(*mock.ResponseWriter)

	if c.Status != 500 {
		t.Fatal("status should be 500")
	}

	if string(w.MockBytes) != "error" {
		t.Fatal("content should be error")
	}
}

func TestForbidden(t *testing.T) {
	c := mock.NewClient()
	Forbidden(c, "forbidden")
	w := c.Writer.(*mock.ResponseWriter)

	if c.Status != 403 {
		t.Fatal("status should be 403")
	}

	if string(w.MockBytes) != "forbidden" {
		t.Fatal("content should be forbidden")
	}
}

func TestTooManyRequests(t *testing.T) {
	c := mock.NewClient()
	TooManyRequests(c, "too many requests")
	w := c.Writer.(*mock.ResponseWriter)

	if c.Status != 429 {
		t.Fatal("status should be 429")
	}

	if string(w.MockBytes) != "too many requests" {
		t.Fatal("content should be too many requests")
	}
}
