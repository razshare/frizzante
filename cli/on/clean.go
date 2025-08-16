package on

import (
	"github.com/razshare/frizzante/cli"
	"github.com/razshare/frizzante/cli/path"
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"os/exec"
	"path/filepath"
)

func Clean(c *cli.Cli) error {
	gobin, err := path.Go(c, ".")
	if err != nil {
		return err
	}

	clean := exec.Command(gobin, "clean")
	clean.Env = append(os.Environ())
	clean.Stderr = os.Stderr
	clean.Stdout = os.Stdout
	clean.Stdin = os.Stdin
	err = clean.Run()
	if err != nil {
		return err
	}

	err = os.RemoveAll(filepath.Join(*c.Flags.App, "dist"))
	if err != nil {
		return err
	}

	err = os.RemoveAll(filepath.Join(*c.Flags.App, "node_modules"))
	if err != nil {
		return err
	}

	err = os.RemoveAll(filepath.Join(".gen", "tmp"))
	if err != nil {
		return err
	}

	err = os.RemoveAll(".vite")
	if err != nil {
		return err
	}

	err = Touch(c)
	if err != nil {
		return err
	}

	messages.Success("project cleaned")

	return nil
}
