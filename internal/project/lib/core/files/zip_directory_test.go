package files

import (
	"os"
	"testing"
)

func TestZipDirectory(t *testing.T) {
	_ = os.Remove("dir.zip")
	defer func() { _ = os.Remove("dir.zip") }()
	err := ZipDirectory("dir", "dir.zip")
	if err != nil {
		t.Fatal(err)
	}
	if !IsFile("dir.zip") {
		t.Fatal("dir.zip should be a file")
	}
}
