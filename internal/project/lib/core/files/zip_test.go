package files

import (
	"os"
	"testing"
)

func TestZipFile(t *testing.T) {
	_ = os.Remove("zip_test.zip")
	defer func() { _ = os.Remove("zip_test.zip") }()

	err := ZipFile("zip_test.go", "zip_test.zip")
	if err != nil {
		t.Fatal(err)
	}

	if !IsFile("zip_test.zip") {
		t.Fatal("zip_test.zip should be a file")
	}

}

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
