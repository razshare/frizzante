package files

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyFile(t *testing.T) {
	_ = os.RemoveAll("test_copy_file")
	defer func() { _ = os.RemoveAll("test_copy_file") }()

	err := CopyFile("copy_test.go", filepath.Join("test_copy_file", "copy_test.go"))
	if err != nil {
		t.Fatal(err)
	}

	if !IsFile(filepath.Join("test_copy_file", "copy_test.go")) {
		t.Fatalf("test_copy_file/copy_test.go should be a file")
	}
}

func TestCopyDirectory(t *testing.T) {
	_ = os.RemoveAll("test_copy_directory")
	defer func() { _ = os.RemoveAll("test_copy_directory") }()

	err := CopyDirectory("dir", "test_copy_directory")
	if err != nil {
		t.Fatal(err)
	}

	if !IsDirectory("test_copy_directory") {
		t.Fatalf("test_copy_directory should be a directory")
	}
}
