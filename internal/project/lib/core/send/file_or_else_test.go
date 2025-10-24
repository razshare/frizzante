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
var EfsTestFileOrElse embed.FS

func TestFileOrElse(t *testing.T) {
	_ = os.Rename("test.txt", "test.renamed.txt")
	defer func() { _ = os.Rename("test.renamed.txt", "test.txt") }()
	client := mocks.NewClient()
	client.Config.Efs = EfsTestFileOrElse
	client.Config.PublicRoot = ""
	client.Request.RequestURI = "test.txt"
	var orElse bool
	FileOrElse(client, func() { orElse = true })
	writer := client.Writer.(*mocks.ResponseWriter)

	if orElse {
		t.Fatal("else branch should not trigger")
	}

	if !strings.Contains(string(writer.MockBytes), "this is a test") {
		t.Fatal("content should contain this is a test")
	}
}

func TestFileOrElseFromFs(t *testing.T) {
	client := mocks.NewClient()
	client.Config.PublicRoot = ""
	client.Request.RequestURI = "test.txt"
	client.Request.URL = &url.URL{Path: "test.txt"}
	var orElse bool
	FileOrElse(client, func() { orElse = true })
	writer := client.Writer.(*mocks.ResponseWriter)

	if orElse {
		t.Fatal("else branch should not trigger")
	}

	if !strings.Contains(string(writer.MockBytes), "this is a test") {
		t.Fatal("content should contain this is a test")
	}
}

func TestFileOrElseShouldFail(t *testing.T) {
	client := mocks.NewClient()
	client.Config.Efs = EfsTestFileOrElse
	client.Config.PublicRoot = ""
	client.Request.RequestURI = "some_file.go"
	var orElse bool
	FileOrElse(client, func() { orElse = true })
	if !orElse {
		t.Fatal("or else should trigger")
	}
}
