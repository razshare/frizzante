package send

import "testing"

func TestHeader(t *testing.T) {
	client := MockClient()
	Header(client, "key", "value")
	writer := client.Writer.(*MockWriter)
	if writer.MockHeader.Get("key") != "value" {
		t.Fatal("key should be value")
	}
}

func TestHeaders(t *testing.T) {
	client := MockClient()
	Headers(client, map[string]string{
		"key1": "value1",
		"key2": "value2",
	})

	writer := client.Writer.(*MockWriter)

	if writer.MockHeader.Get("key1") != "value1" {
		t.Fatal("key1 should be value1")
	}

	if writer.MockHeader.Get("key2") != "value2" {
		t.Fatal("key2 should be value2")
	}
}

func TestRedirect(t *testing.T) {
	client := MockClient()
	Redirect(client, "/about", 303)
	writer := client.Writer.(*MockWriter)

	if client.Status != 303 {
		t.Fatal("status should be 303")
	}

	loc := writer.MockHeader.Get("Location")

	if loc != "/about" {
		t.Fatal("location should be about")
	}
}

func TestNavigate(t *testing.T) {
	client := MockClient()
	Navigate(client, "/about")
	writer := client.Writer.(*MockWriter)

	if client.Status != 302 {
		t.Fatal("status should be 302")
	}

	loc := writer.MockHeader.Get("Location")

	if loc != "/about" {
		t.Fatal("location should be about")
	}
}

func TestContentType(t *testing.T) {
	client := MockClient()
	ContentType(client, "text/html")
	writer := client.Writer.(*MockWriter)

	ctype := writer.MockHeader.Get("Content-Type")

	if ctype != "text/html" {
		t.Fatal("content type should be text/html")
	}
}
