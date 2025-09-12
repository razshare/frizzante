package generate

import (
	"embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/platform"
)

//go:embed internal/project/**
var TestDatabaseEfs embed.FS

func TestDatabase(t *testing.T) {
	if err := os.RemoveAll("lib"); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll("lib") }()

	if err := Database(DatabaseOptions{
		Platform: platform.LinuxAmd64,
		Efs:      TestDatabaseEfs,
		Auto:     true,
		Go:       "go",
		Type:     "sqlite",
	}); err != nil {
		t.Fatal(err)
	}

	if !files.IsFile(filepath.Join("lib", "database", "sqlite", "example.txt")) {
		t.Fatal("lib/database/sqlite/example.txt should exist")
	}

	if !files.IsFile(filepath.Join("lib", "database", "sqlite", "sqlc.yaml")) {
		t.Fatal("lib/database/sqlite/sqlc.yaml should exist")
	}

	if !files.IsFile(filepath.Join("lib", "database", "sqlite", "queries.sql")) {
		t.Fatal("lib/database/sqlite/queries.sql should exist")
	}

	if !files.IsFile(filepath.Join("lib", "database", "sqlite", "schema.sql")) {
		t.Fatal("lib/database/sqlite/schema.sql should exist")
	}
}
