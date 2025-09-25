package send

import (
	"embed"
	"net/url"
	"strings"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/mock"
)

//go:embed file_or_else_test.go
var EfsTestFileOrElse embed.FS

func TestFileOrElse(t *testing.T) {
	client := mock.NewClient()
	client.Config.Efs = EfsTestFileOrElse
	client.Config.PublicRoot = ""
	client.Request.RequestURI = "file_or_else_test.go"
	var orElse bool
	FileOrElse(client, FileOrElseConfig{
		UseDisk: false,
		OrElse:  func() { orElse = true },
	})
	writer := client.Writer.(*mock.ResponseWriter)

	if orElse {
		t.Fatal("else branch should not trigger")
	}

	if !strings.Contains(string(writer.MockBytes), "var EfsTestFileOrElse embed.FS") {
		t.Fatal("content should contain this file")
	}
}

func TestFileOrElseFromFs(t *testing.T) {
	client := mock.NewClient()
	client.Config.PublicRoot = ""
	client.Request.RequestURI = "file_or_else_test.go"
	client.Request.URL = &url.URL{Path: "file_or_else_test.go"}
	var orElse bool
	FileOrElse(client, FileOrElseConfig{
		UseDisk: true,
		OrElse:  func() { orElse = true },
	})
	writer := client.Writer.(*mock.ResponseWriter)

	if orElse {
		t.Fatal("else branch should not trigger")
	}

	if !strings.Contains(string(writer.MockBytes), "var TestFileOrElseFromFs embed.FS") {
		t.Fatal("content should contain this file")
	}
}

func TestFileOrElseShouldFail(t *testing.T) {
	client := mock.NewClient()
	client.Config.Efs = EfsTestFileOrElse
	client.Config.PublicRoot = ""
	client.Request.RequestURI = "some_file.go"
	var orElse bool
	FileOrElse(client, FileOrElseConfig{
		UseDisk: false,
		OrElse:  func() { orElse = true },
	})
	if !orElse {
		t.Fatal("or else should trigger")
	}
}
