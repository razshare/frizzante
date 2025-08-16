package path

import (
	"github.com/razshare/frizzante/cli"
	"github.com/razshare/frizzante/cli/extension"
	"os"
	"path/filepath"
	"strings"
)

func Bun(c *cli.Cli, base string) (string, error) {
	var bin string

	if *c.Flags.Bun != "" {
		bin = *c.Flags.Bun
	} else {
		bin = filepath.Join(".gen", "bun", "bun")
	}

	if strings.HasPrefix(bin, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		bin = strings.Replace(bin, "~", home, 1)
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
