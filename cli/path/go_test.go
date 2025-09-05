package path

import (
	"os"
	"strings"
	"testing"
)

func TestGo(t *testing.T) {
	_go, err := Go("./go/go")
	if err != nil {
		t.Fatal(err)
	}

	if _go != "./go/go" {
		t.Fatal("binary should be ./go/go")
	}
}

func TestGoAtHome(t *testing.T) {
	_go, err := Go("~/.go/go")
	if err != nil {
		t.Fatal(err)
	}

	user, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasPrefix(_go, user) {
		t.Fatal("binary should be prefixed with user home dir")
	}
}
