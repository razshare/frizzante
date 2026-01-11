package generations

import (
	"embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
)

//go:embed internal/additions/**
var TestDatabaseEfs embed.FS

func TestDatabase(t *testing.T) {
	if err := os.RemoveAll(".gen"); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(".gen") }()

	if err := Sqlc(SqlcOptions{Sqlc: filepath.Join(".gen", "sqlc", "sqlc")}); err != nil {
		t.Fatal(err)
	}

	if !files.IsFile(filepath.Join(".gen", "sqlc", "sqlc")) {
		t.Fatal(".gen/sqlc/sqlc should exist")
	}

	if err := os.RemoveAll("lib"); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll("lib") }()

	if err := Databases(DatabasesOptions{
		Efs:  TestDatabaseEfs,
		Go:   "go",
		Type: "sqlite",
		Sqlc: filepath.Join(".gen", "sqlc", "sqlc"),
	}); err != nil {
		t.Fatal(err)
	}

	if !files.IsFile(filepath.Join("lib", "databases", "sqlite", "example.txt")) {
		t.Fatal("lib/databases/sqlite/example.txt should exist")
	}

	if !files.IsFile(filepath.Join("lib", "databases", "sqlite", "queries.sql")) {
		t.Fatal("lib/databases/sqlite/queries.sql should exist")
	}

	if !files.IsFile(filepath.Join("lib", "databases", "sqlite", "schema.sql")) {
		t.Fatal("lib/databases/sqlite/schema.sql should exist")
	}

	if !files.IsFile(filepath.Join("lib", "databases", "sqlite", "sqlc.yaml")) {
		t.Fatal("lib/databases/sqlite/sqlc.yaml should exist")
	}
}
