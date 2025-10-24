package generate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/platforms"
)

func TestBun(t *testing.T) {
	if err := os.RemoveAll(".gen"); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(".gen") }()

	if err := Bun(BunOptions{Platform: platforms.LinuxAmd64, Auto: true, Bun: filepath.Join(".gen", "bun", "bun")}); err != nil {
		t.Fatal(err)
	}

	if !files.IsFile(filepath.Join(".gen", "bun", "bun")) {
		t.Fatal(".gen/bun/bun should exist")
	}
}
