package embeds

import (
	"embed"
	"os"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
)

//go:embed zip_file_test.go
var TestZipFileEfs embed.FS

func TestZipFile(t *testing.T) {
	_ = os.RemoveAll("zip_file_test.zip")
	defer func() { _ = os.RemoveAll("zip_file_test.zip") }()
	if err := ZipFile(TestZipFileEfs, "zip_file_test.go", "zip_file_test.zip"); err != nil {
		t.Fatal(err)
	}
	if !files.IsFile("zip_file_test.zip") {
		t.Fatal("zip_file_test.zip should be a file")
	}
}
