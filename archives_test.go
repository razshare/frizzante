package main

import (
	"github.com/razshare/frizzante/archives"
	"github.com/razshare/frizzante/files"
	"os"
	"path/filepath"
	"testing"
)

func TestSet(t *testing.T) {
	files.DeleteFile(filepath.Join(".gen", "archive", "domain", "key"))
	_ = archives.New(filepath.Join(".gen", "archive")).Set("domain", "key", []byte("content"))
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
	archive := archives.New(filepath.Join(".gen", "archive"))
	_ = archive.Set("domain", "key", []byte("content"))
	readBytes, _ := archive.Get("domain", "key")
	if "content" != string(readBytes) {
		t.Fatal("wrong archive value")
	}
}

func TestHas(t *testing.T) {
	files.DeleteFile(filepath.Join(".gen", "archive", "domain", "key"))
	archive := archives.New(filepath.Join(".gen", "archive"))
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
	archive := archives.New(filepath.Join(".gen", "archive"))
	_ = archive.Set("domain", "key", []byte("content"))
	_ = archive.Remove("domain", "key")
	if files.IsFile("archive/domain/key") {
		t.Fatal("archive should not contain value")
	}
}

func TestHasDomain(t *testing.T) {
	_ = os.RemoveAll(filepath.Join(".gen", "archive", "domain", "key"))
	archive := archives.New(filepath.Join(".gen", "archive"))
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
	archive := archives.New(filepath.Join(".gen", "archive"))
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
