package embeds

import (
	"embed"
	"main/lib/core/files"
	"os"
	"testing"
)

//go:embed zip_test.go
var TestZipFileEfs embed.FS

func TestZipFile(t *testing.T) {
	_ = os.Remove("zip_test.zip")
	defer func() { _ = os.Remove("zip_test.zip") }()

	err := ZipFile(TestZipFileEfs, "zip_test.go", "zip_test.zip")
	if err != nil {
		t.Fatal(err)
	}

	if !files.IsFile("zip_test.zip") {
		t.Fatal("zip_test.zip should be a file")
	}
}

//go:embed dir
var TestZipDirectoryEfs embed.FS

func TestZipDirectory(t *testing.T) {
	_ = os.Remove("dir.zip")
	defer func() { _ = os.Remove("dir.zip") }()

	err := ZipDirectory(TestZipDirectoryEfs, "dir", "dir.zip")
	if err != nil {
		t.Fatal(err)
	}

	if !files.IsFile("dir.zip") {
		t.Fatal("dir.zip should be a file")
	}
}
