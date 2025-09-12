package generate

import (
	"embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/razshare/frizzante/files"
)

//go:embed internal/project/**
var TestFormsEfs embed.FS

func TestForms(t *testing.T) {
	if err := os.RemoveAll("app"); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll("app") }()

	if err := Forms(FormsOptions{Auto: true, Efs: TestFormsEfs, App: "app"}); err != nil {
		t.Fatal(err)
	}

	if !files.IsFile(filepath.Join("app", "lib", "components", "forms", "example.txt")) {
		t.Fatal("app/lib/components/forms/example.txt should exist")
	}
}
