package send

import (
	"embed"
	"net/url"
	"strings"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mocks"
)

//go:embed app
var EfsTestFileOrElse embed.FS

func TestFileOrElse(t *testing.T) {
	client := mocks.NewClient()
	client.Config.Efs = EfsTestFileOrElse
	client.Request.RequestURI = "index.html"
	client.Request.URL = &url.URL{Path: "index.html"}
	var orElse bool
	FileOrElse(client, func() { orElse = true })
	writer := client.Writer.(*mocks.ResponseWriter)

	if orElse {
		t.Fatal("else branch should not trigger")
	}

	if !strings.Contains(string(writer.MockBytes), "<html") {
		t.Fatal("index.html file should contain <html")
	}
}

func TestFileOrElseFromFs(t *testing.T) {
	client := mocks.NewClient()
	client.Request.RequestURI = "index.html"
	client.Request.URL = &url.URL{Path: "index.html"}
	var orElse bool
	FileOrElse(client, func() { orElse = true })
	writer := client.Writer.(*mocks.ResponseWriter)

	if orElse {
		t.Fatal("else branch should not trigger")
	}

	if !strings.Contains(string(writer.MockBytes), "<html") {
		t.Fatal("index.html file should contain <html")
	}
}

func TestFileOrElseShouldFail(t *testing.T) {
	client := mocks.NewClient()
	client.Config.Efs = EfsTestFileOrElse
	client.Request.RequestURI = "some_file.go"
	client.Request.URL = &url.URL{Path: "some_file.go"}
	var orElse bool
	FileOrElse(client, func() { orElse = true })
	if !orElse {
		t.Fatal("or else should trigger")
	}
}
