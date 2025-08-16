package on

import (
	"github.com/razshare/frizzante/cli"
	"github.com/razshare/frizzante/cli/path"
	"os"
	"os/exec"
)

func Check(c *cli.Cli) error {
	err := Touch(c)
	if err != nil {
		return err
	}

	bunbin, err := path.Bun(c, *c.Flags.App)
	if err != nil {
		return err
	}

	eslint := exec.Command(bunbin, "x", "eslint")
	eslint.Dir = *c.Flags.App
	eslint.Env = append(os.Environ())
	eslint.Stderr = os.Stderr
	eslint.Stdout = os.Stdout
	eslint.Stdin = os.Stdin
	err = eslint.Run()
	if err != nil {
		return err
	}

	svelteCheck := exec.Command(bunbin, "x", "svelte-check", "--tsconfig=./tsconfig.json")
	svelteCheck.Dir = *c.Flags.App
	svelteCheck.Env = append(os.Environ())
	svelteCheck.Stderr = os.Stderr
	svelteCheck.Stdout = os.Stdout
	svelteCheck.Stdin = os.Stdin
	err = svelteCheck.Run()
	if err != nil {
		return err
	}
	return nil
}
