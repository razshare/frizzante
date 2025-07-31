package main

import (
	"github.com/razshare/frizzante/archives"
	"github.com/razshare/frizzante/files"
	"os"
	"path/filepath"
	"testing"
)

func TestNewDiskArchive(test *testing.T) {
	_ = os.RemoveAll(filepath.Join(".gen", "archive"))
	archive := archives.NewDiskArchive(filepath.Join(".gen", "archive"))
	if filepath.Join(".gen", "archive") != archive.Name {
		test.Fatalf("archive name must be `%s`", filepath.Join(".gen", "archive"))
	}
	if len(archive.Lock.Names) > 0 {
		test.Fatal("archive lanes must be empty")
	}
	if files.IsDirectory("archive") {
		test.Fatal("archive directory must be created on first set, not on archive creation")
	}
}

func TestSet(test *testing.T) {
	files.DeleteFile(filepath.Join(".gen", "archive", "domain", "key"))
	_ = archives.NewDiskArchive(filepath.Join(".gen", "archive")).Set("domain", "key", []byte("content"))
	if !files.IsFile(filepath.Join(".gen", "archive", "domain", "key")) {
		test.Fatal("domain value is missing")
	}
	readBytes, _ := os.ReadFile(filepath.Join(".gen", "archive", "domain", "key"))
	if "content" != string(readBytes) {
		test.Fatal("wrong archive value")
	}
}

func TestGet(test *testing.T) {
	files.DeleteFile(filepath.Join(".gen", "archive", "domain", "key"))
	archive := archives.NewDiskArchive(filepath.Join(".gen", "archive"))
	_ = archive.Set("domain", "key", []byte("content"))
	readBytes, _ := archive.Get("domain", "key")
	if "content" != string(readBytes) {
		test.Fatal("wrong archive value")
	}
}

func TestHas(test *testing.T) {
	files.DeleteFile(filepath.Join(".gen", "archive", "domain", "key"))
	archive := archives.NewDiskArchive(filepath.Join(".gen", "archive"))
	has, _ := archive.Has("domain", "key")
	if has {
		test.Fatal("archive should not contain value yet")
	}

	_ = archive.Set("domain", "key", []byte("content"))

	has, _ = archive.Has("domain", "key")
	if !has {
		test.Fatal("archive is missing value")
	}
}

func TestRemove(test *testing.T) {
	files.DeleteFile(filepath.Join(".gen", "archive", "domain", "key"))
	archive := archives.NewDiskArchive(filepath.Join(".gen", "archive"))
	_ = archive.Set("domain", "key", []byte("content"))
	_ = archive.Remove("domain", "key")
	if files.IsFile("archive/domain/key") {
		test.Fatal("archive should not contain value")
	}
}

func TestHasDomain(test *testing.T) {
	_ = os.RemoveAll(filepath.Join(".gen", "archive", "domain", "key"))
	archive := archives.NewDiskArchive(filepath.Join(".gen", "archive"))
	_ = archive.Set("domain", "key", []byte("content"))
	has, _ := archive.HasDomain("domain")
	if !has {
		test.Fatal("archive should have domain")
	}
	_ = archive.Remove("domain", "key")
	has, _ = archive.HasDomain("domain")
	if !has {
		test.Fatal("archive should have domain")
	}
}

func TestRemoveDomain(test *testing.T) {
	_ = os.RemoveAll(filepath.Join(".gen", "archive", "domain", "key"))
	archive := archives.NewDiskArchive(filepath.Join(".gen", "archive"))
	_ = archive.Set("domain", "key", []byte("content"))
	_ = archive.RemoveDomain("domain")
	has, _ := archive.Has("domain", "key")
	if has {
		test.Fatal("archive should not have value")
	}
	has, _ = archive.HasDomain("domain")
	if has {
		test.Fatal("archive should not have domain")
	}
}
