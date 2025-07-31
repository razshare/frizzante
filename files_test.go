package main

import (
	"github.com/razshare/frizzante/embeds"
	"github.com/razshare/frizzante/files"
	"testing"
)

func TestEmbeddedIsFile(test *testing.T) {
	// Positive.
	fileName := "makefile"
	actual := embeds.IsFile(efs, fileName)
	expected := true
	if actual != expected {
		test.Fatalf("%s (embedded) was expected to be a file", fileName)
	}

	// Negatives.
	fileName = "t"
	actual = embeds.IsFile(efs, fileName)
	expected = false
	if actual != expected {
		test.Fatalf("%s (embedded) was expected to not be a file", fileName)
	}

	fileName = "qwerty"
	actual = embeds.IsFile(efs, fileName)
	expected = false
	if actual != expected {
		test.Fatalf("%s (embedded) was expected to not be a file", fileName)
	}
}

func TestEmbeddedIsDirectory(test *testing.T) {
	// Positive.
	fileName := ".github"
	actual := embeds.IsDirectory(efs, fileName)
	expected := true
	if actual != expected {
		test.Fatalf("%s (embedded) was expected to be a directory", fileName)
	}

	// Negatives.
	fileName = "makefile"
	actual = embeds.IsDirectory(efs, fileName)
	expected = false
	if actual != expected {
		test.Fatalf("%s (embedded) was expected to not be a directory", fileName)
	}

	fileName = "qwerty"
	actual = embeds.IsDirectory(efs, fileName)
	expected = false
	if actual != expected {
		test.Fatalf("%s (embedded) was expected to not be a directory", fileName)
	}
}

func TestIsFile(test *testing.T) {
	// Positive.
	fileName := "makefile"
	actual := files.IsFile(fileName)
	expected := true
	if actual != expected {
		test.Fatalf("%s was expected to be a file", fileName)
	}

	// Negatives.
	fileName = ".github"
	actual = files.IsFile(fileName)
	expected = false
	if actual != expected {
		test.Fatalf("%s was expected to not be a file", fileName)
	}

	fileName = "qwerty"
	actual = files.IsFile(fileName)
	expected = false
	if actual != expected {
		test.Fatalf("%s was expected to not be a file", fileName)
	}
}

func TestIsDirectory(test *testing.T) {
	// Positive.
	fileName := ".github"
	actual := files.IsDirectory(fileName)
	expected := true
	if actual != expected {
		test.Fatalf("%s was expected to be a directory", fileName)
	}

	// Negatives.
	fileName = "makefile"
	actual = files.IsDirectory(fileName)
	expected = false
	if actual != expected {
		test.Fatalf("%s was expected to not be a directory", fileName)
	}

	fileName = "qwerty"
	actual = files.IsDirectory(fileName)
	expected = false
	if actual != expected {
		test.Fatalf("%s was expected to not be a directory", fileName)
	}
}
