package embeds

import (
	"embed"
	"testing"
)

//go:embed dir
//go:embed is_file_test.go
var TestIsDirectoryEfs embed.FS

func TestIsDirectory(t *testing.T) {
	if !IsDirectory(TestIsDirectoryEfs, "dir") {
		t.Fatal("dir should be a directory")
	}
	if IsDirectory(TestIsDirectoryEfs, "is_file_test.go") {
		t.Fatal("is_file_test.go should not be a directory")
	}
	if IsDirectory(TestIsDirectoryEfs, "qwerty") {
		t.Fatal("qwerty should not be a directory")
	}
}
