package main

import (
	"github.com/razshare/frizzante/archives"
	"github.com/razshare/frizzante/files"
	"os"
	"path/filepath"
	"testing"
)

func TestNewDiskArchive(t *testing.T) {
	_ = os.RemoveAll(filepath.Join(".gen", "archive"))
	archive := archives.NewDiskArchive()
	if filepath.Join(".gen", "archive") != archive.Name {
		t.Fatalf("archive name must be `%s`", filepath.Join(".gen", "archive"))
	}
	if len(archive.Road.Lanes) > 0 {
		t.Fatal("archive lanes must be empty")
	}
	if files.IsDirectory("archive") {
		t.Fatal("archive directory must be created on first set, not on archive creation")
	}
}

func TestSet(t *testing.T) {
	files.DeleteFile(filepath.Join(".gen", "archive", "domain", "key"))
	_ = archives.NewDiskArchive().Set("domain", "key", []byte("content"))
	if !files.IsFile(filepath.Join(".gen", "archive", "domain", "key")) {
		t.Fatal("domain value is missing")
	}
	readBytes, _ := os.ReadFile(filepath.Join(".gen", "archive", "domain", "key"))
	if "content" != string(readBytes) {
		t.Fatal("wrong archive value")
	}
}

func TestGet(t *testing.T) {
	files.DeleteFile(filepath.Join(".gen", "archive", "domain", "key"))
	archive := archives.NewDiskArchive()
	_ = archive.Set("domain", "key", []byte("content"))
	readBytes, _ := archive.Get("domain", "key")
	if "content" != string(readBytes) {
		t.Fatal("wrong archive value")
	}
}

func TestHas(t *testing.T) {
	files.DeleteFile(filepath.Join(".gen", "archive", "domain", "key"))
	archive := archives.NewDiskArchive()
	has, _ := archive.Has("domain", "key")
	if has {
		t.Fatal("archive should not contain value yet")
	}

	_ = archive.Set("domain", "key", []byte("content"))

	has, _ = archive.Has("domain", "key")
	if !has {
		t.Fatal("archive is missing value")
	}
}

func TestRemove(t *testing.T) {
	files.DeleteFile(filepath.Join(".gen", "archive", "domain", "key"))
	archive := archives.NewDiskArchive()
	_ = archive.Set("domain", "key", []byte("content"))
	_ = archive.Remove("domain", "key")
	if files.IsFile("archive/domain/key") {
		t.Fatal("archive should not contain value")
	}
}

func TestHasDomain(t *testing.T) {
	_ = os.RemoveAll(filepath.Join(".gen", "archive", "domain", "key"))
	archive := archives.NewDiskArchive()
	_ = archive.Set("domain", "key", []byte("content"))
	has, _ := archive.HasDomain("domain")
	if !has {
		t.Fatal("archive should have domain")
	}
	_ = archive.Remove("domain", "key")
	has, _ = archive.HasDomain("domain")
	if !has {
		t.Fatal("archive should have domain")
	}
}

func TestRemoveDomain(t *testing.T) {
	_ = os.RemoveAll(filepath.Join(".gen", "archive", "domain", "key"))
	archive := archives.NewDiskArchive()
	_ = archive.Set("domain", "key", []byte("content"))
	_ = archive.RemoveDomain("domain")
	has, _ := archive.Has("domain", "key")
	if has {
		t.Fatal("archive should not have value")
	}
	has, _ = archive.HasDomain("domain")
	if has {
		t.Fatal("archive should not have domain")
	}
}
