package main

import (
	"github.com/razshare/frizzante/embeds"
	"github.com/razshare/frizzante/files"
	"testing"
)

func TestEmbeddedIsFile(t *testing.T) {
	// Positive.
	n := "makefile"
	ac := embeds.IsFile(Tefs, n)
	ex := true
	if ac != ex {
		t.Fatalf("%s (embedded) was expected to be a file", n)
	}

	// Negatives.
	n = "t"
	ac = embeds.IsFile(Tefs, n)
	ex = false
	if ac != ex {
		t.Fatalf("%s (embedded) was expected to not be a file", n)
	}

	n = "qwerty"
	ac = embeds.IsFile(Tefs, n)
	ex = false
	if ac != ex {
		t.Fatalf("%s (embedded) was expected to not be a file", n)
	}
}

func TestEmbeddedIsDirectory(t *testing.T) {
	// Positive.
	n := ".github"
	ac := embeds.IsDirectory(Tefs, n)
	ex := true
	if ac != ex {
		t.Fatalf("%s (embedded) was expected to be a directory", n)
	}

	// Negatives.
	n = "makefile"
	ac = embeds.IsDirectory(Tefs, n)
	ex = false
	if ac != ex {
		t.Fatalf("%s (embedded) was expected to not be a directory", n)
	}

	n = "qwerty"
	ac = embeds.IsDirectory(Tefs, n)
	ex = false
	if ac != ex {
		t.Fatalf("%s (embedded) was expected to not be a directory", n)
	}
}

func TestIsFile(t *testing.T) {
	// Positive.
	n := "makefile"
	ac := files.IsFile(n)
	ex := true
	if ac != ex {
		t.Fatalf("%s was expected to be a file", n)
	}

	// Negatives.
	n = ".github"
	ac = files.IsFile(n)
	ex = false
	if ac != ex {
		t.Fatalf("%s was expected to not be a file", n)
	}

	n = "qwerty"
	ac = files.IsFile(n)
	ex = false
	if ac != ex {
		t.Fatalf("%s was expected to not be a file", n)
	}
}

func TestIsDirectory(t *testing.T) {
	// Positive.
	n := ".github"
	ac := files.IsDirectory(n)
	ex := true
	if ac != ex {
		t.Fatalf("%s was expected to be a directory", n)
	}

	// Negatives.
	n = "makefile"
	ac = files.IsDirectory(n)
	ex = false
	if ac != ex {
		t.Fatalf("%s was expected to not be a directory", n)
	}

	n = "qwerty"
	ac = files.IsDirectory(n)
	ex = false
	if ac != ex {
		t.Fatalf("%s was expected to not be a directory", n)
	}
}
