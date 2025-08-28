package files

import (
	"os"
	"path/filepath"
	"testing"
)

func TestUnzipFile(t *testing.T) {
	_ = os.RemoveAll("test_unzip_file_dir")
	defer func() { _ = os.RemoveAll("test_unzip_file_dir") }()

	if err := UnzipFile("unzip_test.zip", "test_unzip_file_dir"); err != nil {
		t.Fatal(err)
	}

	if !IsFile(filepath.Join("test_unzip_file_dir", "unzip_test.go")) {
		t.Fatal("test_unzip_file_dir/unzip_test.go should be a file")
	}
}

func TestUnzipFileShouldFail(t *testing.T) {
	_ = os.RemoveAll("test_unzip_file_dir")
	defer func() { _ = os.RemoveAll("test_unzip_file_dir") }()

	err := os.MkdirAll("test_unzip_file_dir", os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	if err = os.WriteFile(filepath.Join("test_unzip_file_dir", "dir"), []byte("hello"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	if err = UnzipFile("unzip_test.zip", "test_unzip_file_dir"); err == nil {
		t.Fatal("unzip should fail")
	}
}
