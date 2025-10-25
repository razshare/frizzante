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
var TestQueriesEfs embed.FS

func TestQueries(t *testing.T) {
	if err := os.RemoveAll(".gen"); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll("lib"); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(".gen") }()
	defer func() { _ = os.RemoveAll("lib") }()

	if err := Sqlc(SqlcOptions{
		Auto:     true,
		Sqlc:     filepath.Join(".gen", "sqlc", "sqlc"),
		Platform: platforms.LinuxAmd64,
	}); err != nil {
		t.Fatal(err)
	}

	if !files.IsFile(filepath.Join(".gen", "sqlc", "sqlc")) {
		t.Fatal(".gen/sqlc/sqlc should exist")
	}

	if err := Database(DatabaseOptions{
		Platform: platforms.LinuxAmd64,
		Efs:      TestQueriesEfs,
		Auto:     true,
		Go:       "go",
		Type:     "sqlite",
		Sqlc:     filepath.Join(".gen", "sqlc", "sqlc"),
	}); err != nil {
		t.Fatal(err)
	}

	if err := Queries(QueriesOptions{
		Auto:     true,
		Sqlc:     filepath.Join(".gen", "sqlc", "sqlc"),
		Platform: platforms.LinuxAmd64,
	}); err != nil {
		t.Fatal(err)
	}

	if !files.IsFile(filepath.Join(".gen", "sqlc", "sqlc")) {
		t.Fatal(".gen/sqlc/sqlc should exist")
	}

	if !files.IsFile(filepath.Join("lib", "sqlite", "databases", "sqlc", "db.go")) {
		t.Fatal("lib/sqlite/databases/sqlc/db.go should exist")
	}

	if !files.IsFile(filepath.Join("lib", "sqlite", "databases", "sqlc", "models.go")) {
		t.Fatal("lib/sqlite/databases/sqlc/models.go should exist")
	}

	if !files.IsFile(filepath.Join("lib", "sqlite", "databases", "sqlc", "queries.sql.go")) {
		t.Fatal("lib/sqlite/databases/sqlc/queries.sql.go should exist")
	}
}
