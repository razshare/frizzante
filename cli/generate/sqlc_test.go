package generate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/platform"
)

func TestSqlc(t *testing.T) {
	if err := os.RemoveAll(".gen"); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(".gen") }()

	if err := Sqlc(SqlcOptions{
		Auto:     true,
		Sqlc:     filepath.Join(".gen", "sqlc", "sqlc"),
		Platform: platform.LinuxAmd64,
	}); err != nil {
		t.Fatal(err)
	}

	if !files.IsFile(filepath.Join(".gen", "sqlc", "sqlc")) {
		t.Fatal(".gen/sqlc/sqlc should exist")
	}
}
