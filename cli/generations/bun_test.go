package generations

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/razshare/frizzante/v2/internal/project/lib/core/files"
)

func TestBun(t *testing.T) {
	if err := os.RemoveAll(".gen"); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(".gen") }()
	if err := Bun(BunOptions{Bun: filepath.Join(".gen", "bun", "bun")}); err != nil {
		t.Fatal(err)
	}
	if !files.IsFile(filepath.Join(".gen", "bun", "bun")) {
		t.Fatal(".gen/bun/bun should exist")
	}
}
