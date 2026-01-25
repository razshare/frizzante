package embeds

import (
	"embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
)

//go:embed dir
//go:embed copy_file_test.go
var TestCopyFileEfs embed.FS

func TestCopyFile(t *testing.T) {
	_ = os.RemoveAll("test_copy_file_dir")
	defer func() { _ = os.RemoveAll("test_copy_file_dir") }()

	err := CopyFile(TestCopyFileEfs, "copy_file_test.go", filepath.Join("test_copy_file_dir", "copy_file_test.go"))
	if err != nil {
		t.Fatal(err)
	}

	if !files.IsFile(filepath.Join("test_copy_file_dir", "copy_file_test.go")) {
		t.Fatalf("test_copy_file_dir/copy_file_test.go should be a file")
	}

	err = CopyFile(TestCopyFileEfs, "dir/test.txt", filepath.Join("test_copy_file_dir", "copy_file_test.go"))
	if err != nil {
		t.Fatal(err)
	}

	d, err := os.ReadFile(filepath.Join("test_copy_file_dir", "copy_file_test.go"))
	if err != nil {
		t.Fatal(err)
	}

	if string(d) != "hello world" {
		t.Fatal("test_copy_file_dir/copy_file_test.go should contain hello world")
	}
}
