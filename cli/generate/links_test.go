package generate

import (
	"embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
)

//go:embed internal/additions/**
var TestLinksEfs embed.FS

func TestLinks(t *testing.T) {
	if err := os.RemoveAll("app"); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll("app") }()

	if err := Links(LinksOptions{Auto: true, Efs: TestLinksEfs, App: "app"}); err != nil {
		t.Fatal(err)
	}

	if !files.IsFile(filepath.Join("app", "lib", "components", "links", "example.txt")) {
		t.Fatal("app/lib/components/links/example.txt should exist")
	}
}
