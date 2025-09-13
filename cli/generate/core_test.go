package generate

import (
	"embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
)

//go:embed internal/project/**
var TestCoreEfs embed.FS

func TestCore(t *testing.T) {
	if err := os.RemoveAll("lib"); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll("app"); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll("lib") }()
	defer func() { _ = os.RemoveAll("app") }()

	if err := Core(CoreOptions{App: "app", Auto: true, Efs: TestCoreEfs}); err != nil {
		t.Fatal(err)
	}

	if !files.IsFile(filepath.Join("lib", "core", "server", "example.txt")) {
		t.Fatal("lib/core/server/example.txt should exist")
	}

	if !files.IsFile(filepath.Join("app", "lib", "scripts", "core", "example.txt")) {
		t.Fatal("app/lib/scripts/core/example.txt should exist")
	}
}
