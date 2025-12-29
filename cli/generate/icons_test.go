package generate

import (
	"embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
)

//go:embed internal/project/**
var TestIconsEfs embed.FS

func TestIcons(t *testing.T) {
	var err error
	if err = os.RemoveAll("app"); err != nil {
		t.Fatal(err)
	}
	if err = os.RemoveAll(".gen"); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll("app") }()
	defer func() { _ = os.RemoveAll(".gen") }()

	if err = Bun(BunOptions{Bun: filepath.Join(".gen", "bun", "bun")}); err != nil {
		return
	}

	if err = Icons(IconsOptions{
		Efs: TestIconsEfs,
		Bun: filepath.Join(".gen", "bun", "bun"),
	}); err != nil {
		t.Fatal(err)
	}

	if !files.IsFile(filepath.Join("app", "lib", "components", "icons", "example.txt")) {
		t.Fatal("app/lib/components/icons/example.txt should exist")
	}
}
