package send

import (
	"main/lib/core/mock"
	"testing"
)

func TestHeader(t *testing.T) {
	client := mock.NewClient()
	Header(client, "key", "value")
	writer := client.Writer.(*mock.ResponseWriter)
	if writer.MockHeader.Get("key") != "value" {
		t.Fatal("key should be value")
	}
}

func TestHeaders(t *testing.T) {
	client := mock.NewClient()
	Headers(client, map[string]string{
		"key1": "value1",
		"key2": "value2",
	})

	writer := client.Writer.(*mock.ResponseWriter)

	if writer.MockHeader.Get("key1") != "value1" {
		t.Fatal("key1 should be value1")
	}

	if writer.MockHeader.Get("key2") != "value2" {
		t.Fatal("key2 should be value2")
	}
}

func TestRedirect(t *testing.T) {
	client := mock.NewClient()
	Redirect(client, "/about", 303)
	writer := client.Writer.(*mock.ResponseWriter)

	if client.Status != 303 {
		t.Fatal("status should be 303")
	}

	if writer.MockHeader.Get("Location") != "/about" {
		t.Fatal("location should be about")
	}
}

func TestNavigate(t *testing.T) {
	client := mock.NewClient()
	Navigate(client, "/about")
	writer := client.Writer.(*mock.ResponseWriter)

	if client.Status != 302 {
		t.Fatal("status should be 302")
	}

	if writer.MockHeader.Get("Location") != "/about" {
		t.Fatal("location should be about")
	}
}

func TestContentType(t *testing.T) {
	client := mock.NewClient()
	ContentType(client, "text/html")
	writer := client.Writer.(*mock.ResponseWriter)

	if writer.MockHeader.Get("Content-Type") != "text/html" {
		t.Fatal("content type should be text/html")
	}
}
