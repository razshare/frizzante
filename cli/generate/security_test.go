package generate

import (
	"embed"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
)

//go:embed internal/additions/**
//go:embed internal/project/**
var TestSecurityEfs embed.FS

func TestSecurity(t *testing.T) {
	var err error
	if err = os.RemoveAll("lib"); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll("lib") }()

	if err = Security(SecurityOptions{
		Efs:  TestSecurityEfs,
		Auto: true,
	}); err != nil {
		t.Fatal(err)
	}

	if !files.IsFile(filepath.Join("lib", "security", "example.go")) {
		t.Fatal("lib/security/example.go should exist")
	}

	var data []byte
	if data, err = os.ReadFile(filepath.Join("lib", "security", "example.go")); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(data), "main/lib/core/files") {
		t.Fatal("example.go should contain main/lib/core/files")
	}
}
