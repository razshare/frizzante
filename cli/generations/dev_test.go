package generations

import (
	"embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
)

//go:embed internal/project/**
var TestDevEfs embed.FS

func TestDev(t *testing.T) {
	if err := os.RemoveAll("lib"); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll("app"); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll("lib") }()
	defer func() { _ = os.RemoveAll("app") }()
	if err := Dev(DevOptions{Efs: TestDevEfs}); err != nil {
		t.Fatal(err)
	}
	if !files.IsFile(filepath.Join("lib", "dev", "types", "example.txt")) {
		t.Fatal("lib/dev/types/example.txt should exist")
	}
}
