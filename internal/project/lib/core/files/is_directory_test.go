package files

import "testing"

func TestIsDirectory(t *testing.T) {
	if !IsDirectory("dir") {
		t.Fatal("dir should be a directory")
	}
	if IsDirectory("files_test.go") {
		t.Fatal("files_test.go should not be a directory")
	}
	if IsDirectory("qwerty") {
		t.Fatal("qwerty should not be a directory")
	}
}
