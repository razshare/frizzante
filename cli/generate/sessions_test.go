package generate

import (
	"embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
)

//go:embed internal/additions/**
//go:embed internal/project/**
var TestSessionEfs embed.FS

func TestSession(t *testing.T) {
	if err := os.RemoveAll("lib"); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll("lib") }()

	if err := Sessions(SessionsOptions{
		Efs:  TestSessionEfs,
		Auto: true,
		Type: "memory",
	}); err != nil {
		t.Fatal(err)
	}

	if !files.IsFile(filepath.Join("lib", "memory", "sessions", "example.txt")) {
		t.Fatal("lib/memory/sessions/example.txt should exist")
	}

	if err := Sessions(SessionsOptions{
		Efs:  TestSessionEfs,
		Auto: true,
		Type: "disk",
	}); err != nil {
		t.Fatal(err)
	}

	if !files.IsFile(filepath.Join("lib", "disk", "sessions", "example.txt")) {
		t.Fatal("lib/disk/sessions/example.txt should exist")
	}
}
