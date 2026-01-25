package files

import (
	"os"
	"testing"
)

func TestCopyDirectory(t *testing.T) {
	_ = os.RemoveAll("test_copy_file_dir")
	defer func() { _ = os.RemoveAll("test_copy_file_dir") }()
	err := CopyDirectory("dir", "test_copy_file_dir")
	if err != nil {
		t.Fatal(err)
	}
	if !IsDirectory("test_copy_file_dir") {
		t.Fatalf("test_copy_file_dir should be a directory")
	}
}
