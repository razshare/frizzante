package generate

import (
	"embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
)

//go:embed internal/project/**
var TestProjectEfs embed.FS

func TestProject(t *testing.T) {
	if err := os.RemoveAll("asd"); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll("asd") }()

	if err := Project(ProjectOptions{Auto: true, Efs: TestProjectEfs, Name: "asd"}); err != nil {
		t.Fatal(err)
	}

	if !files.IsFile(filepath.Join("asd", "lib", "core", "server", "example.txt")) {
		t.Fatal("asd/lib/core/server/example.txt should exist")
	}

	if !files.IsFile(filepath.Join("asd", "app", "lib", "scripts", "core", "example.txt")) {
		t.Fatal("asd/app/lib/scripts/core/example.txt should exist")
	}

	if files.IsFile(filepath.Join("asd", "lib", "session", "disk", "example.txt")) {
		t.Fatal("asd/lib/session/disk/example.txt should not exist")
	}

	if !files.IsFile(filepath.Join("asd", "lib", "session", "memory", "example.txt")) {
		t.Fatal("asd/lib/session/memory/example.txt should not exist")
	}

	if files.IsFile(filepath.Join("asd", "lib", "database", "sqlite", "example.txt")) {
		t.Fatal("asd/lib/database/sqlite/example.txt should not exist")
	}

	if files.IsFile(filepath.Join("asd", "app", "lib", "components", "links", "example.txt")) {
		t.Fatal("asd/app/lib/components/links/example.txt should not exist")
	}

	if files.IsFile(filepath.Join("asd", "app", "lib", "components", "forms", "example.txt")) {
		t.Fatal("asd/app/lib/components/forms/example.txt should not exist")
	}
}
