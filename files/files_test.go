package files

import "testing"

func TestIsFile(t *testing.T) {
	// Positive.
	n := "files_test.go"
	ac := IsFile(n)
	ex := true
	if ac != ex {
		t.Fatalf("%s was expected to be a file", n)
	}

	// Negatives.
	n = "dir"
	ac = IsFile(n)
	ex = false
	if ac != ex {
		t.Fatalf("%s was expected to not be a file", n)
	}

	n = "qwerty"
	ac = IsFile(n)
	ex = false
	if ac != ex {
		t.Fatalf("%s was expected to not be a file", n)
	}
}

func TestIsDirectory(t *testing.T) {
	// Positive.
	n := "dir"
	ac := IsDirectory(n)
	ex := true
	if ac != ex {
		t.Fatalf("%s was expected to be a directory", n)
	}

	// Negatives.
	n = "files_test.go"
	ac = IsDirectory(n)
	ex = false
	if ac != ex {
		t.Fatalf("%s was expected to not be a directory", n)
	}

	n = "qwerty"
	ac = IsDirectory(n)
	ex = false
	if ac != ex {
		t.Fatalf("%s was expected to not be a directory", n)
	}
}
