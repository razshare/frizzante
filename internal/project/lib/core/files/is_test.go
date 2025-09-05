package files

import "testing"

func TestIsFile(t *testing.T) {
	if !IsFile("files_test.go") != true {
		t.Fatal("files_test.go should be a file")
	}

	if IsFile("dir") {
		t.Fatal("dir should not be a file")
	}

	if IsFile("qwerty") {
		t.Fatal("qwerty should not be a file")
	}
}

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
