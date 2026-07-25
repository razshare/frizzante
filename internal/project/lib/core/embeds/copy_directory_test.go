package embeds

import (
	"embed"
	"os"
	"testing"

	"github.com/razshare/frizzante/v2/internal/project/lib/core/files"
)

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
