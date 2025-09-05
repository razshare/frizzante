package embeds

import (
	"embed"
	"main/lib/core/files"
	"os"
	"path/filepath"
	"testing"
)

//go:embed dir
//go:embed copy_test.go
var TestCopyFileEfs embed.FS

func TestCopyFile(t *testing.T) {
	_ = os.RemoveAll("test_copy_file_dir")
	defer func() { _ = os.RemoveAll("test_copy_file_dir") }()

	err := CopyFile(TestCopyFileEfs, "copy_test.go", filepath.Join("test_copy_file_dir", "copy_test.go"))
	if err != nil {
		t.Fatal(err)
	}

	if !files.IsFile(filepath.Join("test_copy_file_dir", "copy_test.go")) {
		t.Fatalf("test_copy_file_dir/copy_test.go should be a file")
	}

	err = CopyFile(TestCopyFileEfs, filepath.Join("dir", "test.txt"), filepath.Join("test_copy_file_dir", "copy_test.go"))
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

//go:embed dir
var TestCopyDirectoryEfs embed.FS

func TestCopyDirectory(t *testing.T) {
	_ = os.RemoveAll("test_copy_file_dir")
	defer func() { _ = os.RemoveAll("test_copy_file_dir") }()

	err := CopyDirectory(TestCopyDirectoryEfs, "dir", "test_copy_file_dir")
	if err != nil {
		t.Fatal(err)
	}

	if !files.IsDirectory("test_copy_file_dir") {
		t.Fatalf("test_copy_file_dir should be a directory")
	}
}
