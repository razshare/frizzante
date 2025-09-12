package generate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/platform"
)

func TestQueries(t *testing.T) {
	if err := os.RemoveAll(".gen"); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll("lib"); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(".gen") }()
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

	if err := Queries(QueriesOptions{
		Auto:     true,
		Sqlc:     filepath.Join(".gen", "sqlc", "sqlc"),
		Platform: platform.LinuxAmd64,
	}); err != nil {
		t.Fatal(err)
	}

	if !files.IsFile(filepath.Join(".gen", "sqlc", "sqlc")) {
		t.Fatal(".gen/sqlc/sqlc should exist")
	}

	if !files.IsFile(filepath.Join("lib", "database", "sqlite", "sqlc", "db.go")) {
		t.Fatal("lib/database/sqlite/sqlc/db.go should exist")
	}

	if !files.IsFile(filepath.Join("lib", "database", "sqlite", "sqlc", "models.go")) {
		t.Fatal("lib/database/sqlite/sqlc/models.go should exist")
	}

	if !files.IsFile(filepath.Join("lib", "database", "sqlite", "sqlc", "queries.sql.go")) {
		t.Fatal("lib/database/sqlite/sqlc/queries.sql.go should exist")
	}
}
