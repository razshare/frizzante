package frizzante

import (
	"github.com/razshare/frizzante/libarchive"
	"github.com/razshare/frizzante/libfs"
	"os"
	"testing"
)

func TestNewDiskArchive(t *testing.T) {
	_ = os.RemoveAll("archive")
	archive := libarchive.NewDiskArchive()
	if "archive" != archive.Name {
		t.Fatal("archive name must be `archive`")
	}
	if len(archive.Road.Lanes) > 0 {
		t.Fatal("archive lanes must be empty")
	}
	if libfs.IsDirectory("archive") {
		t.Fatal("archive directory must be created on first set, not on archive creation")
	}
}

func TestWithName(t *testing.T) {
	archive := libarchive.NewDiskArchive().WithName("test")
	if "test" != archive.Name {
		t.Fatal("archive name must be `test`")
	}
}

func TestSet(t *testing.T) {
	libfs.DeleteFile("archive/domain/key")
	archive := libarchive.NewDiskArchive()
	_ = archive.Set("domain", "key", []byte("content"))
	if !libfs.IsFile("archive/domain/key") {
		t.Fatal("domain value is missing")
	}
	readBytes, _ := os.ReadFile("archive/domain/key")
	if "content" != string(readBytes) {
		t.Fatal("wrong archive value")
	}
}

func TestGet(t *testing.T) {
	libfs.DeleteFile("archive/domain/key")
	archive := libarchive.NewDiskArchive()
	_ = archive.Set("domain", "key", []byte("content"))
	readBytes, _ := archive.Get("domain", "key")
	if "content" != string(readBytes) {
		t.Fatal("wrong archive value")
	}
}

func TestHas(t *testing.T) {
	libfs.DeleteFile("archive/domain/key")
	archive := libarchive.NewDiskArchive()
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
	libfs.DeleteFile("archive/domain/key")
	archive := libarchive.NewDiskArchive()
	_ = archive.Set("domain", "key", []byte("content"))
	_ = archive.Remove("domain", "key")
	if libfs.IsFile("archive/domain/key") {
		t.Fatal("archive should not contain value")
	}
}

func TestHasDomain(t *testing.T) {
	_ = os.RemoveAll("archive/domain/key")
	archive := libarchive.NewDiskArchive()
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
	_ = os.RemoveAll("archive/domain/key")
	archive := libarchive.NewDiskArchive()
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
