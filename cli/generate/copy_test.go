package generate

import (
	"embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/razshare/frizzante/files"
)

//go:embed air.go
var TestCopyEfs embed.FS

func TestCopy(t *testing.T) {
	if err := os.RemoveAll(".gen"); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(".gen") }()

	if err := Copy(CopyOptions{From: "air.go", To: filepath.Join(".gen", "air.go.txt"), Efs: TestCopyEfs}); err != nil {
		t.Fatal(err)
	}

	if !files.IsFile(filepath.Join(".gen", "air.go.txt")) {
		t.Fatal(".gen/air.go.txt should exist")
	}
}
