package send

import (
	"errors"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

func TestContent(t *testing.T) {
	client := mocks.NewClient()
	Content(client, []byte("hello"))
	writer := client.Writer.(*mocks.ResponseWriter)

	if string(writer.MockBytes) != "hello" {
		t.Fatal("content should be hello")
	}
}

func TestMessage(t *testing.T) {
	client := mocks.NewClient()
	Message(client, "hello")
	writer := client.Writer.(*mocks.ResponseWriter)

	if string(writer.MockBytes) != "hello" {
		t.Fatal("content should be hello")
	}
}

func TestMessagef(t *testing.T) {
	client := mocks.NewClient()
	Messagef(client, "hello %s", "world")
	writer := client.Writer.(*mocks.ResponseWriter)

	if string(writer.MockBytes) != "hello world" {
		t.Fatal("content should be hello world")
	}
}

func TestNotFound(t *testing.T) {
	client := mocks.NewClient()
	NotFound(client, "not found")
	writer := client.Writer.(*mocks.ResponseWriter)

	if client.Status != 404 {
		t.Fatal("status should be 404")
	}

	if string(writer.MockBytes) != "not found" {
		t.Fatal("content should be not found")
	}
}

func TestUnauthorized(t *testing.T) {
	client := mocks.NewClient()
	Unauthorized(client, "unauthorized")
	writer := client.Writer.(*mocks.ResponseWriter)

	if client.Status != 401 {
		t.Fatal("status should be 401")
	}

	if string(writer.MockBytes) != "unauthorized" {
		t.Fatal("content should be unauthorized")
	}
}

func TestBadRequest(t *testing.T) {
	client := mocks.NewClient()
	BadRequest(client, "bad request")
	writer := client.Writer.(*mocks.ResponseWriter)

	if client.Status != 400 {
		t.Fatal("status should be 400")
	}

	if string(writer.MockBytes) != "bad request" {
		t.Fatal("content should be bad request")
	}
}

func TestError(t *testing.T) {
	client := mocks.NewClient()
	Error(client, errors.New("error"))
	writer := client.Writer.(*mocks.ResponseWriter)

	if client.Status != 500 {
		t.Fatal("status should be 500")
	}

	if string(writer.MockBytes) != "error" {
		t.Fatal("content should be error")
	}
}

func TestForbidden(t *testing.T) {
	client := mocks.NewClient()
	Forbidden(client, "forbidden")
	writer := client.Writer.(*mocks.ResponseWriter)

	if client.Status != 403 {
		t.Fatal("status should be 403")
	}

	if string(writer.MockBytes) != "forbidden" {
		t.Fatal("content should be forbidden")
	}
}

func TestTooManyRequests(t *testing.T) {
	client := mocks.NewClient()
	TooManyRequests(client, "too many requests")
	writer := client.Writer.(*mocks.ResponseWriter)

	if client.Status != 429 {
		t.Fatal("status should be 429")
	}

	if string(writer.MockBytes) != "too many requests" {
		t.Fatal("content should be too many requests")
	}
}
