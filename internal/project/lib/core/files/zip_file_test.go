package files

import (
	"os"
	"testing"
)

func TestZipFile(t *testing.T) {
	_ = os.RemoveAll("zip_file_test.zip")
	defer func() { _ = os.RemoveAll("zip_file_test.zip") }()
	err := ZipFile("zip_file_test.go", "zip_file_test.zip")
	if err != nil {
		t.Fatal(err)
	}
	if !IsFile("zip_file_test.zip") {
		t.Fatal("zip_file_test.zip should be a file")
	}
}
