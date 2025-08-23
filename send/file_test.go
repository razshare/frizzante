package send

import (
	"embed"
	"strings"
	"testing"
)

//go:embed file_test.go
var EfsTestFileOrElse embed.FS

func TestFileOrElse(t *testing.T) {
	client := MockClient()
	client.Config.Efs = EfsTestFileOrElse
	// we don't have an app directory in this package,
	// so we need to ignore the public root.
	client.Config.PublicRoot = ""
	// we're intentionally omitting the leading "/",
	// otherwise the embedded file system will not
	// find the "requested" file.
	client.Request.RequestURI = "file_test.go"
	var or bool
	FileOrElse(client, func() { or = true })
	writer := client.Writer.(*MockWriter)

	if or {
		t.Fatal("else branch should not trigger")
	}

	if !strings.Contains(string(writer.MockBytes), "var EfsTestFileOrElse embed.FS") {
		t.Fatal("content should contain this file")
	}
}

func TestFileOrElseShouldFail(t *testing.T) {
	client := MockClient()
	client.Config.Efs = EfsTestFileOrElse
	client.Config.PublicRoot = ""
	client.Request.RequestURI = "some_file.go"
	var or bool
	FileOrElse(client, func() { or = true })
	if !or {
		t.Fatal("or else should trigger")
	}
}
