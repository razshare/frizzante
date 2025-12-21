package apps

import (
	"path/filepath"
	"testing"

	"github.com/razshare/frizzante/cli/extensions"
)

func TestNew(t *testing.T) {
	app := New()

	if *app.Go != "go"+extensions.Find() {
		t.Fatalf("go should be go%s", extensions.Find())
	}

	if *app.Air != filepath.Join(".gen", "air", "air"+extensions.Find()) {
		t.Fatalf("air should be .gen/air/air%s", extensions.Find())
	}

	if *app.Bun != filepath.Join(".gen", "bun", "bun"+extensions.Find()) {
		t.Fatalf("bun should be .gen/air/air%s", extensions.Find())
	}

	if *app.Sqlc != filepath.Join(".gen", "sqlc", "sqlc"+extensions.Find()) {
		t.Fatalf("sqlc should be ./gen/sqlc/sqlc%s", extensions.Find())
	}

	if *app.Tags != "" {
		t.Fatalf("tags should be empty")
	}

	if *app.DatabaseConnectionString != "" {
		t.Fatal("database should be empty")
	}

	if *app.DatabaseType != "sqlc" {
		t.Fatal("database type should be sqlc")
	}
}
