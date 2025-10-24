package send

import (
	"embed"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

//go:embed test.txt
var EfsRequestedFile embed.FS

func TestRequestedFile(t *testing.T) {
	_ = os.Rename("test.txt", "test.renamed.txt")
	defer func() { _ = os.Rename("test.renamed.txt", "test.txt") }()
	client := mocks.NewClient()
	client.Config.Efs = EfsRequestedFile
	client.Config.PublicRoot = ""
	client.Request.RequestURI = "test.txt"
	writer := client.Writer.(*mocks.ResponseWriter)

	if !RequestedFile(client) {
		t.Fatal("sending file should fail succeed")
	}

	if !strings.Contains(string(writer.MockBytes), "this is a test") {
		t.Fatal("content should contain this is a test")
	}
}

func TestRequestedFileFromFs(t *testing.T) {
	client := mocks.NewClient()
	client.Config.PublicRoot = ""
	client.Request.RequestURI = "test.txt"
	client.Request.URL = &url.URL{Path: "test.txt"}
	writer := client.Writer.(*mocks.ResponseWriter)

	if !RequestedFile(client) {
		t.Fatal("sending file should fail succeed")
	}

	if !strings.Contains(string(writer.MockBytes), "this is a test") {
		t.Fatal("content should contain this is a test")
	}
}

func TestRequestedFileShouldFail(t *testing.T) {
	client := mocks.NewClient()
	client.Config.Efs = EfsRequestedFile
	client.Config.PublicRoot = ""
	client.Request.RequestURI = "some_file.go"
	if RequestedFile(client) {
		t.Fatal("sending file should fail")
	}
}
