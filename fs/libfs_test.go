package fs

import (
	"embed"
	"testing"
)

//go:embed mimes.go
//go:embed test/*
var Efs embed.FS

func TestEmbeddedExists(test *testing.T) {
	// Positive.
	fileName := "mimes.go"
	actual := EfsFileExists(Efs, fileName)
	expected := true
	if !actual {
		test.Fatalf("%s (embedded) was expected to exist", fileName)
	}

	// Negative.
	fileName = "qwerty"
	actual = EfsFileExists(Efs, fileName)
	expected = false
	if actual != expected {
		test.Fatalf("%s (embedded) was expected to exist", fileName)
	}
}

func TestEmbeddedIsFile(test *testing.T) {
	// Positive.
	fileName := "mimes.go"
	actual := EfsIsFile(Efs, fileName)
	expected := true
	if actual != expected {
		test.Fatalf("%s (embedded) was expected to be a file", fileName)
	}

	// Negatives.
	fileName = "test"
	actual = EfsIsFile(Efs, fileName)
	expected = false
	if actual != expected {
		test.Fatalf("%s (embedded) was expected to not be a file", fileName)
	}

	fileName = "qwerty"
	actual = EfsIsFile(Efs, fileName)
	expected = false
	if actual != expected {
		test.Fatalf("%s (embedded) was expected to not be a file", fileName)
	}
}

func TestEmbeddedIsDirectory(test *testing.T) {
	// Positive.
	fileName := "test"
	actual := EfsIsDirectory(Efs, fileName)
	expected := true
	if actual != expected {
		test.Fatalf("%s (embedded) was expected to be a directory", fileName)
	}

	// Negatives.
	fileName = "mimes.go"
	actual = EfsIsDirectory(Efs, fileName)
	expected = false
	if actual != expected {
		test.Fatalf("%s (embedded) was expected to not be a directory", fileName)
	}

	fileName = "qwerty"
	actual = EfsIsDirectory(Efs, fileName)
	expected = false
	if actual != expected {
		test.Fatalf("%s (embedded) was expected to not be a directory", fileName)
	}
}

func TestExists(test *testing.T) {
	// Positive.
	fileName := "mimes.go"
	actual := FileExists(fileName)
	expected := true
	if actual != expected {
		test.Fatalf("%s was expected to exist", fileName)
	}

	// Negative.
	fileName = "qwerty"
	actual = FileExists(fileName)
	expected = false
	if actual != expected {
		test.Fatalf("%s was expected to not exist", fileName)
	}
}

func TestIsFile(test *testing.T) {
	// Positive.
	fileName := "mimes.go"
	actual := IsFile(fileName)
	expected := true
	if actual != expected {
		test.Fatalf("%s was expected to be a file", fileName)
	}

	// Negatives.
	fileName = "test"
	actual = IsFile(fileName)
	expected = false
	if actual != expected {
		test.Fatalf("%s was expected to not be a file", fileName)
	}

	fileName = "qwerty"
	actual = IsFile(fileName)
	expected = false
	if actual != expected {
		test.Fatalf("%s was expected to not be a file", fileName)
	}
}

func TestIsDirectory(test *testing.T) {
	// Positive.
	fileName := "test"
	actual := IsDirectory(fileName)
	expected := true
	if actual != expected {
		test.Fatalf("%s was expected to be a directory", fileName)
	}

	// Negatives.
	fileName = "mimes.go"
	actual = IsDirectory(fileName)
	expected = false
	if actual != expected {
		test.Fatalf("%s was expected to not be a directory", fileName)
	}

	fileName = "qwerty"
	actual = IsDirectory(fileName)
	expected = false
	if actual != expected {
		test.Fatalf("%s was expected to not be a directory", fileName)
	}
}
