//go:build queries

package databases

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/sqlc-dev/sqlc/pkg/cli"
)

func Generate() (err error) {
	if code := cli.Run([]string{
		fmt.Sprintf("--file=\"%s\"", filepath.Join("lib", "sqlite", "databases", "sqlc.yaml")),
		"generate",
	}); code != 0 {
		err = errors.New("could not generate queries")
	}
	return
}
