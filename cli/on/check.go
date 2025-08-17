package on

import (
	"os"
	"os/exec"
)

func Check(app string, bunbin string) error {
	err := Touch(app)
	if err != nil {
		return err
	}

	eslint := exec.Command(bunbin, "x", "eslint")
	eslint.Dir = app
	eslint.Env = append(os.Environ())
	eslint.Stderr = os.Stderr
	eslint.Stdout = os.Stdout
	eslint.Stdin = os.Stdin
	err = eslint.Run()
	if err != nil {
		return err
	}

	svelteCheck := exec.Command(bunbin, "x", "svelte-check", "--tsconfig=./tsconfig.json")
	svelteCheck.Dir = app
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
