package on

import (
	"github.com/razshare/frizzante/cli/path"
	"github.com/razshare/frizzante/cli/state"
	"os"
	"os/exec"
)

func Check() error {
	err := Touch()
	if err != nil {
		return err
	}

	bunbin, err := path.Bun(*state.App)
	if err != nil {
		return err
	}

	eslint := exec.Command(bunbin, "x", "eslint")
	eslint.Dir = *state.App
	eslint.Env = append(os.Environ())
	eslint.Stderr = os.Stderr
	eslint.Stdout = os.Stdout
	eslint.Stdin = os.Stdin
	err = eslint.Run()
	if err != nil {
		return err
	}

	svelteCheck := exec.Command(bunbin, "x", "svelte-check", "--tsconfig=./tsconfig.json")
	svelteCheck.Dir = *state.App
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
