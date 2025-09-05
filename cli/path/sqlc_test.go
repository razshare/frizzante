package path

import (
	"os"
	"strings"
	"testing"
)

func TestSqlc(t *testing.T) {
	sqlc, err := Sqlc("./sqlc/sqlc")
	if err != nil {
		t.Fatal(err)
	}

	if sqlc != "./sqlc/sqlc" {
		t.Fatal("binary should be ./sqlc/sqlc")
	}
}

func TestSqlcAtHome(t *testing.T) {
	sqlc, err := Sqlc("~/.sqlc/sqlc")
	if err != nil {
		t.Fatal(err)
	}

	user, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}

	if !strings.HasPrefix(sqlc, user) {
		t.Fatal("binary should be prefixed with user home dir")
	}
}
