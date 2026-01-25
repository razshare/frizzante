package generations

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
)

func TestSqlc(t *testing.T) {
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
}
