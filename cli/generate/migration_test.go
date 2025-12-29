package generate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
)

func TestMigration(t *testing.T) {
	additionsDirectory := filepath.Join("internal", "additions")
	databasesDirectory := filepath.Join(additionsDirectory, "lib", "databases", "sqlite")
	migrationsDirectory := filepath.Join(databasesDirectory, "migrations")
	sqlcYamlFile := filepath.Join(databasesDirectory, "sqlc.yaml")
	sqlcFile := filepath.Join(additionsDirectory, ".gen", "sqlc", "sqlc")

	if err := os.RemoveAll(migrationsDirectory); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := os.RemoveAll(migrationsDirectory); err != nil {
			t.Fatal(err)
		}
	}()

	if err := Migration(MigrationOptions{
		Sqlc:     sqlcFile,
		SqlcYaml: sqlcYamlFile,
	}); err != nil {
		t.Fatal(err)
	}

	if !files.IsDirectory(migrationsDirectory) {
		t.Fatalf("directory %s should exist", migrationsDirectory)
	}
}
