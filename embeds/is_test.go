package embeds

import (
	"embed"
	"testing"
)

//go:embed is_test.go
var TestIsFileEfs embed.FS

func TestIsFile(t *testing.T) {
	if !IsFile(TestIsFileEfs, "is_test.go") {
		t.Fatal("is_test.go should be a file")
	}

	if IsFile(TestIsFileEfs, "qwerty") {
		t.Fatal("qwerty should not be a file")
	}
}

//go:embed dir
//go:embed is_test.go
var TestIsDirectoryEfs embed.FS

func TestIsDirectory(t *testing.T) {
	if !IsDirectory(TestIsDirectoryEfs, "dir") {
		t.Fatal("dir should be a directory")
	}

	if IsDirectory(TestIsDirectoryEfs, "is_test.go") {
		t.Fatal("is_test.go should not be a directory")
	}

	if IsDirectory(TestIsDirectoryEfs, "qwerty") {
		t.Fatal("qwerty should not be a directory")
	}
}
