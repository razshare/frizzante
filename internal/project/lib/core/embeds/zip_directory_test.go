package embeds

import (
	"embed"
	"os"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
)

//go:embed dir
var TestZipDirectoryEfs embed.FS

func TestZipDirectory(t *testing.T) {
	_ = os.Remove("dir.zip")
	defer func() { _ = os.Remove("dir.zip") }()
	err := ZipDirectory(TestZipDirectoryEfs, "dir", "dir.zip")
	if err != nil {
		t.Fatal(err)
	}
	if !files.IsFile("dir.zip") {
		t.Fatal("dir.zip should be a file")
	}
}
