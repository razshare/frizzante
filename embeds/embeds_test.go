package embeds

import (
	"embed"
	"github.com/razshare/frizzante/files"
	"testing"
)

//go:embed dir
//go:embed embeds_test.go
var Efs embed.FS

func TestEmbeddedIsFile(t *testing.T) {
	// Positive.
	n := "embeds_test.go"
	ac := IsFile(Efs, n)
	ex := true
	if ac != ex {
		t.Fatalf("%s (embedded) was expected to be a file", n)
	}

	// Negatives.
	n = "t"
	ac = IsFile(Efs, n)
	ex = false
	if ac != ex {
		t.Fatalf("%s (embedded) was expected to not be a file", n)
	}

	n = "qwerty"
	ac = IsFile(Efs, n)
	ex = false
	if ac != ex {
		t.Fatalf("%s (embedded) was expected to not be a file", n)
	}
}

func TestEmbeddedIsDirectory(t *testing.T) {
	// Positive.
	n := "dir"
	ac := IsDirectory(Efs, n)
	ex := true
	if ac != ex {
		t.Fatalf("%s (embedded) was expected to be a directory", n)
	}

	// Negatives.
	n = "embeds_test.go"
	ac = IsDirectory(Efs, n)
	ex = false
	if ac != ex {
		t.Fatalf("%s (embedded) was expected to not be a directory", n)
	}

	n = "qwerty"
	ac = IsDirectory(Efs, n)
	ex = false
	if ac != ex {
		t.Fatalf("%s (embedded) was expected to not be a directory", n)
	}
}

func TestIsFile(t *testing.T) {
	// Positive.
	n := "embeds_test.go"
	ac := files.IsFile(n)
	ex := true
	if ac != ex {
		t.Fatalf("%s was expected to be a file", n)
	}

	// Negatives.
	n = "dir"
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
	n := "dir"
	ac := files.IsDirectory(n)
	ex := true
	if ac != ex {
		t.Fatalf("%s was expected to be a directory", n)
	}

	// Negatives.
	n = "embeds_test.go"
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
