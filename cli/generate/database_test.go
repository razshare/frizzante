package generate

import (
	"embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/platform"
)

//go:embed internal/additions/**
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
		Sqlc:     filepath.Join(".gen", "sqlc", "sqlc"),
	}); err != nil {
		t.Fatal(err)
	}

	if !files.IsFile(filepath.Join("lib", "database", "example.txt")) {
		t.Fatal("lib/database/example.txt should exist")
	}

	if !files.IsFile(filepath.Join("lib", "database", "sqlc.yaml")) {
		t.Fatal("lib/database/sqlc.yaml should exist")
	}

	if !files.IsFile(filepath.Join("lib", "database", "queries.sql")) {
		t.Fatal("lib/database/queries.sql should exist")
	}

	if !files.IsFile(filepath.Join("lib", "database", "schema.sql")) {
		t.Fatal("lib/database/schema.sql should exist")
	}
}
