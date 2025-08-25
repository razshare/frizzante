package embeds

import (
	"embed"
	"github.com/razshare/frizzante/files"
	"os"
	"path/filepath"
	"testing"
)

//go:embed dir
var TestCopyFileEfs embed.FS

func TestCopyFile(t *testing.T) {
	_ = os.RemoveAll("test_copy_file_dir")
	defer func() { _ = os.RemoveAll("test_copy_file_dir") }()

	err := CopyFile(TestCopyFileEfs, "dir/test.txt", filepath.Join("test_copy_file_dir", "test.txt"))
	if err != nil {
		t.Fatal(err)
	}

	if !files.IsFile(filepath.Join("test_copy_file_dir", "test.txt")) {
		t.Fatal("test_copy_file_dir/test.txt should be a file")
	}
}
