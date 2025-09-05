package path

import (
	"os"
	"strings"
	"testing"
)

func TestAir(t *testing.T) {
	air, err := Air("./air/air")
	if err != nil {
		t.Fatal(err)
	}

	if air != "./air/air" {
		t.Fatal("binary should be ./air/air")
	}
}

func TestAirAtHome(t *testing.T) {
	air, err := Air("~/.air/air")
	if err != nil {
		t.Fatal(err)
	}

	user, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasPrefix(air, user) {
		t.Fatal("binary should be prefixed with user home dir")
	}
}
