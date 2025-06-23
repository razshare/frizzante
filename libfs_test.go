package frizzante

import (
	"github.com/razshare/frizzante/libfs"
	"testing"
)

func TestEmbeddedIsFile(test *testing.T) {
	// Positive.
	fileName := "makefile"
	actual := libfs.EfsIsFile(efs, fileName)
	expected := true
	if actual != expected {
		test.Fatalf("%s (embedded) was expected to be a file", fileName)
	}

	// Negatives.
	fileName = "test"
	actual = libfs.EfsIsFile(efs, fileName)
	expected = false
	if actual != expected {
		test.Fatalf("%s (embedded) was expected to not be a file", fileName)
	}

	fileName = "qwerty"
	actual = libfs.EfsIsFile(efs, fileName)
	expected = false
	if actual != expected {
		test.Fatalf("%s (embedded) was expected to not be a file", fileName)
	}
}

func TestEmbeddedIsDirectory(test *testing.T) {
	// Positive.
	fileName := "app"
	actual := libfs.EfsIsDirectory(efs, fileName)
	expected := true
	if actual != expected {
		test.Fatalf("%s (embedded) was expected to be a directory", fileName)
	}

	// Negatives.
	fileName = "makefile"
	actual = libfs.EfsIsDirectory(efs, fileName)
	expected = false
	if actual != expected {
		test.Fatalf("%s (embedded) was expected to not be a directory", fileName)
	}

	fileName = "qwerty"
	actual = libfs.EfsIsDirectory(efs, fileName)
	expected = false
	if actual != expected {
		test.Fatalf("%s (embedded) was expected to not be a directory", fileName)
	}
}

func TestIsFile(test *testing.T) {
	// Positive.
	fileName := "makefile"
	actual := libfs.IsFile(fileName)
	expected := true
	if actual != expected {
		test.Fatalf("%s was expected to be a file", fileName)
	}

	// Negatives.
	fileName = "app"
	actual = libfs.IsFile(fileName)
	expected = false
	if actual != expected {
		test.Fatalf("%s was expected to not be a file", fileName)
	}

	fileName = "qwerty"
	actual = libfs.IsFile(fileName)
	expected = false
	if actual != expected {
		test.Fatalf("%s was expected to not be a file", fileName)
	}
}

func TestIsDirectory(test *testing.T) {
	// Positive.
	fileName := "app"
	actual := libfs.IsDirectory(fileName)
	expected := true
	if actual != expected {
		test.Fatalf("%s was expected to be a directory", fileName)
	}

	// Negatives.
	fileName = "makefile"
	actual = libfs.IsDirectory(fileName)
	expected = false
	if actual != expected {
		test.Fatalf("%s was expected to not be a directory", fileName)
	}

	fileName = "qwerty"
	actual = libfs.IsDirectory(fileName)
	expected = false
	if actual != expected {
		test.Fatalf("%s was expected to not be a directory", fileName)
	}
}
