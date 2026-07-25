package generations

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestFixImports(t *testing.T) {
	var err error
	if err = os.RemoveAll(".gen"); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(".gen") }()
	data := []byte("package main\n\nimport \"github.com/razshare/frizzante/v2/internal/project/lib/core/server\"")
	if err = os.MkdirAll(".gen", os.ModePerm); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(".gen", "main.go"), data, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	if err = FixImports(FixImportsOptions{Directory: ".gen"}); err != nil {
		t.Fatal(err)
	}
	if data, err = os.ReadFile(filepath.Join(".gen", "main.go")); err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(data, []byte("github.com/razshare/frizzante/v2/internal/project/lib/core/server")) {
		t.Fatal(".gen/main.go should not reference a global import")
	}
	if !bytes.Contains(data, []byte("main/lib/core/server")) {
		t.Fatal(".gen/main.go should reference a local import")
	}
}
