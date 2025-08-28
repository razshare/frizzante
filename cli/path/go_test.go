package path

import (
	"os"
	"strings"
	"testing"
)

func TestGo(t *testing.T) {
	_go, err := Air("./go/go")
	if err != nil {
		t.Fatal(err)
	}

	if _go != "./go/go" {
		t.Fatal("binary should be ./go/go")
	}
}

func TestGoAtHome(t *testing.T) {
	_go, err := Air("~/.go/go")
	if err != nil {
		t.Fatal(err)
	}

	usr, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasPrefix(_go, usr) {
		t.Fatal("binary should be prefixed with user home dir")
	}
}
