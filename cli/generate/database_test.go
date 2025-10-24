package generate

import (
	"embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/platforms"
)

//go:embed internal/additions/**
var TestDatabaseEfs embed.FS

func TestDatabase(t *testing.T) {
	if err := os.RemoveAll("lib"); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll("lib") }()

	if err := Database(DatabaseOptions{
		Platform: platforms.LinuxAmd64,
		Efs:      TestDatabaseEfs,
		Auto:     true,
		Go:       "go",
		Type:     "sqlite",
		Sqlc:     filepath.Join(".gen", "sqlc", "sqlc"),
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

	if !files.IsFile(filepath.Join("lib", "database", "sqlite", "schema.up.sql")) {
		t.Fatal("lib/database/sqlite/schema.up.sql should exist")
	}
}
