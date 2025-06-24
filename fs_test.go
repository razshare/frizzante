package main

import (
	"github.com/razshare/frizzante/efs"
	"github.com/razshare/frizzante/fs"
	"testing"
)

func TestEmbeddedIsFile(test *testing.T) {
	// Positive.
	fileName := "makefile"
	actual := efs.IsFile(emb, fileName)
	expected := true
	if actual != expected {
		test.Fatalf("%s (embedded) was expected to be a file", fileName)
	}

	// Negatives.
	fileName = "test"
	actual = efs.IsFile(emb, fileName)
	expected = false
	if actual != expected {
		test.Fatalf("%s (embedded) was expected to not be a file", fileName)
	}

	fileName = "qwerty"
	actual = efs.IsFile(emb, fileName)
	expected = false
	if actual != expected {
		test.Fatalf("%s (embedded) was expected to not be a file", fileName)
	}
}

func TestEmbeddedIsDirectory(test *testing.T) {
	// Positive.
	fileName := "app"
	actual := efs.IsDirectory(emb, fileName)
	expected := true
	if actual != expected {
		test.Fatalf("%s (embedded) was expected to be a directory", fileName)
	}

	// Negatives.
	fileName = "makefile"
	actual = efs.IsDirectory(emb, fileName)
	expected = false
	if actual != expected {
		test.Fatalf("%s (embedded) was expected to not be a directory", fileName)
	}

	fileName = "qwerty"
	actual = efs.IsDirectory(emb, fileName)
	expected = false
	if actual != expected {
		test.Fatalf("%s (embedded) was expected to not be a directory", fileName)
	}
}

func TestIsFile(test *testing.T) {
	// Positive.
	fileName := "makefile"
	actual := fs.IsFile(fileName)
	expected := true
	if actual != expected {
		test.Fatalf("%s was expected to be a file", fileName)
	}

	// Negatives.
	fileName = "app"
	actual = fs.IsFile(fileName)
	expected = false
	if actual != expected {
		test.Fatalf("%s was expected to not be a file", fileName)
	}

	fileName = "qwerty"
	actual = fs.IsFile(fileName)
	expected = false
	if actual != expected {
		test.Fatalf("%s was expected to not be a file", fileName)
	}
}

func TestIsDirectory(test *testing.T) {
	// Positive.
	fileName := "app"
	actual := fs.IsDirectory(fileName)
	expected := true
	if actual != expected {
		test.Fatalf("%s was expected to be a directory", fileName)
	}

	// Negatives.
	fileName = "makefile"
	actual = fs.IsDirectory(fileName)
	expected = false
	if actual != expected {
		test.Fatalf("%s was expected to not be a directory", fileName)
	}

	fileName = "qwerty"
	actual = fs.IsDirectory(fileName)
	expected = false
	if actual != expected {
		test.Fatalf("%s was expected to not be a directory", fileName)
	}
}
