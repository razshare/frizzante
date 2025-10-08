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
	if err := os.RemoveAll("app"); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll("app") }()

	if err := Icons(IconsOptions{Auto: true, Efs: TestIconsEfs, App: "app"}); err != nil {
		t.Fatal(err)
	}

	if !files.IsFile(filepath.Join("app", "lib", "components", "icons", "example.txt")) {
		t.Fatal("app/lib/components/icons/example.txt should exist")
	}
}
