package generate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/platforms"
)

func TestMigration(t *testing.T) {
	additionsDirectory := filepath.Join("internal", "additions")
	databasesDirectory := filepath.Join(additionsDirectory, "lib", "sqlite", "databases")
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

	options := MigrationOptions{
		Sqlc:     sqlcFile,
		SqlcYaml: sqlcYamlFile,
		Platform: platforms.Detect(),
	}

	if err := Migration(options); err != nil {
		t.Fatal(err)
	}

	if !files.IsDirectory(migrationsDirectory) {
		t.Fatalf("directory %s should exist", migrationsDirectory)
	}
}
