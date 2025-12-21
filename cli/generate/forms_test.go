package generate

import (
	"embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
)

//go:embed internal/additions/**
var TestFormsEfs embed.FS

func TestForms(t *testing.T) {
	if err := os.RemoveAll("app"); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll("app") }()

	if err := Forms(FormsOptions{
		Efs: TestFormsEfs,
	}); err != nil {
		t.Fatal(err)
	}

	if !files.IsFile(filepath.Join("app", "lib", "components", "forms", "example.txt")) {
		t.Fatal("app/lib/components/forms/example.txt should exist")
	}
}
