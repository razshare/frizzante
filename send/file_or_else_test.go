package send

import (
	"embed"
	"github.com/razshare/frizzante/mock"
	"net/url"
	"strings"
	"testing"
)

//go:embed file_or_else_test.go
var EfsTestFileOrElse embed.FS

func TestFileOrElse(t *testing.T) {
	c := mock.NewClient()
	c.Config.Efs = EfsTestFileOrElse
	c.Config.PublicRoot = ""
	c.Request.RequestURI = "file_or_else_test.go"
	var or bool
	FileOrElse(c, func() { or = true })
	w := c.Writer.(*mock.ResponseWriter)

	if or {
		t.Fatal("else branch should not trigger")
	}

	if !strings.Contains(string(w.MockBytes), "var EfsTestFileOrElse embed.FS") {
		t.Fatal("content should contain this file")
	}
}

func TestFileOrElseFromFs(t *testing.T) {
	c := mock.NewClient()
	c.Config.PublicRoot = ""
	c.Request.RequestURI = "file_or_else_test.go"
	c.Request.URL = &url.URL{Path: "file_or_else_test.go"}
	var or bool
	FileOrElse(c, func() { or = true })
	w := c.Writer.(*mock.ResponseWriter)

	if or {
		t.Fatal("else branch should not trigger")
	}

	if !strings.Contains(string(w.MockBytes), "var TestFileOrElseFromFs embed.FS") {
		t.Fatal("content should contain this file")
	}
}

func TestFileOrElseShouldFail(t *testing.T) {
	c := mock.NewClient()
	c.Config.Efs = EfsTestFileOrElse
	c.Config.PublicRoot = ""
	c.Request.RequestURI = "some_file.go"
	var or bool
	FileOrElse(c, func() { or = true })
	if !or {
		t.Fatal("or else should trigger")
	}
}
