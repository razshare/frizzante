package send

import (
	"embed"
	"net/url"
	"strings"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/dev/mocks"
)

//go:embed app
var EfsRequestedFile embed.FS

func TestRequestedFile(t *testing.T) {
	client := mocks.NewClient()
	client.Options.Efs = EfsRequestedFile
	client.Request.RequestURI = "index.html"
	client.Request.URL = &url.URL{Path: "index.html"}
	writer := client.Writer.(*mocks.ResponseWriter)

	if !RequestedFile(client) {
		t.Fatal("sending file should succeed")
	}

	if !strings.Contains(string(writer.MockBytes), "<html") {
		t.Fatal("index.html file should contain <html")
	}
}

func TestRequestedFileFromFs(t *testing.T) {
	client := mocks.NewClient()
	client.Request.RequestURI = "index.html"
	client.Request.URL = &url.URL{Path: "index.html"}
	writer := client.Writer.(*mocks.ResponseWriter)

	if !RequestedFile(client) {
		t.Fatal("sending file should fail succeed")
	}

	if !strings.Contains(string(writer.MockBytes), "<html") {
		t.Fatal("index.html file should contain <html")
	}
}

func TestRequestedFileShouldFail(t *testing.T) {
	client := mocks.NewClient()
	client.Options.Efs = EfsRequestedFile
	client.Request.RequestURI = "some_file.go"
	client.Request.URL = &url.URL{Path: "some_file.go"}
	if RequestedFile(client) {
		t.Fatal("sending file should fail")
	}
}
