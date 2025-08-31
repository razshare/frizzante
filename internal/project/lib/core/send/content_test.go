package send

import (
	"errors"
	"main/lib/core/mock"
	"testing"
)

func TestContent(t *testing.T) {
	client := mock.NewClient()
	Content(client, []byte("hello"))
	writer := client.Writer.(*mock.ResponseWriter)

	if string(writer.MockBytes) != "hello" {
		t.Fatal("content should be hello")
	}
}

func TestMessage(t *testing.T) {
	client := mock.NewClient()
	Message(client, "hello")
	writer := client.Writer.(*mock.ResponseWriter)

	if string(writer.MockBytes) != "hello" {
		t.Fatal("content should be hello")
	}
}

func TestMessagef(t *testing.T) {
	client := mock.NewClient()
	Messagef(client, "hello %s", "world")
	writer := client.Writer.(*mock.ResponseWriter)

	if string(writer.MockBytes) != "hello world" {
		t.Fatal("content should be hello world")
	}
}

func TestNotFound(t *testing.T) {
	client := mock.NewClient()
	NotFound(client, "not found")
	writer := client.Writer.(*mock.ResponseWriter)

	if client.Status != 404 {
		t.Fatal("status should be 404")
	}

	if string(writer.MockBytes) != "not found" {
		t.Fatal("content should be not found")
	}
}

func TestUnauthorized(t *testing.T) {
	client := mock.NewClient()
	Unauthorized(client, "unauthorized")
	writer := client.Writer.(*mock.ResponseWriter)

	if client.Status != 401 {
		t.Fatal("status should be 401")
	}

	if string(writer.MockBytes) != "unauthorized" {
		t.Fatal("content should be unauthorized")
	}
}

func TestBadRequest(t *testing.T) {
	client := mock.NewClient()
	BadRequest(client, "bad request")
	writer := client.Writer.(*mock.ResponseWriter)

	if client.Status != 400 {
		t.Fatal("status should be 400")
	}

	if string(writer.MockBytes) != "bad request" {
		t.Fatal("content should be bad request")
	}
}

func TestError(t *testing.T) {
	client := mock.NewClient()
	Error(client, errors.New("error"))
	writer := client.Writer.(*mock.ResponseWriter)

	if client.Status != 500 {
		t.Fatal("status should be 500")
	}

	if string(writer.MockBytes) != "error" {
		t.Fatal("content should be error")
	}
}

func TestForbidden(t *testing.T) {
	client := mock.NewClient()
	Forbidden(client, "forbidden")
	writer := client.Writer.(*mock.ResponseWriter)

	if client.Status != 403 {
		t.Fatal("status should be 403")
	}

	if string(writer.MockBytes) != "forbidden" {
		t.Fatal("content should be forbidden")
	}
}

func TestTooManyRequests(t *testing.T) {
	client := mock.NewClient()
	TooManyRequests(client, "too many requests")
	writer := client.Writer.(*mock.ResponseWriter)

	if client.Status != 429 {
		t.Fatal("status should be 429")
	}

	if string(writer.MockBytes) != "too many requests" {
		t.Fatal("content should be too many requests")
	}
}
