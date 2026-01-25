package generations

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
)

func TestAir(t *testing.T) {
	if err := os.RemoveAll(".gen"); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(filepath.Join(".gen")) }()
	if err := Air(AirOptions{Air: filepath.Join(".gen", "air", "air")}); err != nil {
		t.Fatal(err)
	}
	if !files.IsFile(filepath.Join(".gen", "air", "air")) {
		t.Fatal(".gen/air/air should exist")
	}
}
