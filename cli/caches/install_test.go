package caches

import (
	"path/filepath"
	"testing"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
)

func TestInstall(t *testing.T) {
	if err := Install(InstallOptions{
		FromFileName:    "test.zip",
		ToDirectoryName: filepath.Join(".gen", "download"),
	}); err != nil {
		t.Fatal(err)
	}
	if !files.IsDirectory(".gen/download") {
		t.Fatal(".gen/download should exist")
	}
}
