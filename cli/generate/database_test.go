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
	if err := os.RemoveAll(".gen"); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(".gen") }()

	if err := Sqlc(SqlcOptions{
		Sqlc:     filepath.Join(".gen", "sqlc", "sqlc"),
		Platform: platforms.Detect(),
	}); err != nil {
		t.Fatal(err)
	}

	if !files.IsFile(filepath.Join(".gen", "sqlc", "sqlc")) {
		t.Fatal(".gen/sqlc/sqlc should exist")
	}

	if err := os.RemoveAll("lib"); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll("lib") }()

	if err := Database(DatabaseOptions{
		Platform: platforms.Detect(),
		Efs:      TestDatabaseEfs,
		Go:       "go",
		Type:     "sqlite",
		Sqlc:     filepath.Join(".gen", "sqlc", "sqlc"),
	}); err != nil {
		t.Fatal(err)
	}

	if !files.IsFile(filepath.Join("lib", "sqlite", "databases", "example.txt")) {
		t.Fatal("lib/sqlite/databases/example.txt should exist")
	}

	if !files.IsFile(filepath.Join("lib", "sqlite", "databases", "queries.sql")) {
		t.Fatal("lib/sqlite/databases/queries.sql should exist")
	}

	if !files.IsFile(filepath.Join("lib", "sqlite", "databases", "schema.sql")) {
		t.Fatal("lib/sqlite/databases/schema.sql should exist")
	}

	if !files.IsFile(filepath.Join("lib", "sqlite", "databases", "sqlc.yaml")) {
		t.Fatal("lib/sqlite/databases/sqlc.yaml should exist")
	}
}
