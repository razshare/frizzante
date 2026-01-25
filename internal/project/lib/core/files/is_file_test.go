package files

import "testing"

func TestIsFile(t *testing.T) {
	if !IsFile("is_file_test.go") {
		t.Fatal("is_file_test.go should be a file")
	}
	if IsFile("dir") {
		t.Fatal("dir should not be a file")
	}
	if IsFile("qwerty") {
		t.Fatal("qwerty should not be a file")
	}
}
