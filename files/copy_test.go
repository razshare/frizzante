package files

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyFile(t *testing.T) {
	_ = os.RemoveAll("test_copy_file_dir")
	defer func() { _ = os.RemoveAll("test_copy_file_dir") }()

	err := CopyFile("copy_test.go", filepath.Join("test_copy_file_dir", "copy_test.go"))
	if err != nil {
		t.Fatal(err)
	}

	if !IsFile(filepath.Join("test_copy_file_dir", "copy_test.go")) {
		t.Fatalf("test_copy_file_dir/copy_test.go should be a file")
	}

	err = CopyFile(filepath.Join("dir", "test.txt"), filepath.Join("test_copy_file_dir", "copy_test.go"))
	if err != nil {
		t.Fatal(err)
	}

	d, err := os.ReadFile(filepath.Join("test_copy_file_dir", "copy_test.go"))
	if err != nil {
		t.Fatal(err)
	}

	if string(d) != "hello world" {
		t.Fatal("test_copy_file_dir/copy_test.go should contain hello world")
	}
}

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
