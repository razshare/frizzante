package path

import (
	"github.com/razshare/frizzante/cli"
	"github.com/razshare/frizzante/cli/extension"
	"os"
	"path/filepath"
	"strings"
)

func Sqlc(c *cli.Cli, base string) (string, error) {
	var bin string

	if *c.Flags.Sqlc != "" {
		bin = *c.Flags.Sqlc
	} else {
		bin = filepath.Join(".gen", "sqlc", "sqlc")
	}

	if strings.HasPrefix(bin, "~") {
		dirname, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		bin = strings.Replace(bin, "~", dirname, 1)
		return bin + extension.Find(), nil
	}

	if !strings.Contains(bin, string(filepath.Separator)) {
		return bin + extension.Find(), nil
	}

	bin, err := filepath.Rel(base, bin)
	if err != nil {
		return "", err
	}

	return bin + extension.Find(), nil
}
