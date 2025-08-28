package path

import (
	"os"
	"strings"
	"testing"
)

func TestBun(t *testing.T) {
	bun, err := Bun("./bun/bun")
	if err != nil {
		t.Fatal(err)
	}

	if bun != "./bun/bun" {
		t.Fatal("binary should be ./bun/bun")
	}
}

func TestBunAtHome(t *testing.T) {
	bun, err := Bun("~/.bun/bun")
	if err != nil {
		t.Fatal(err)
	}

	usr, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasPrefix(bun, usr) {
		t.Fatal("binary should be prefixed with user home dir")
	}
}
