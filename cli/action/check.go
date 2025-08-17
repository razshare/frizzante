package action

import (
	"os"
	"os/exec"
	"path/filepath"
)

func Check(o CheckOptions) error {
	err := Touch(TouchOptions{App: o.App})
	if err != nil {
		return err
	}

	bun, err := filepath.Rel(o.App, o.Bun)
	if err != nil {
		return err
	}

	eslint := exec.Command(bun, "x", "eslint")
	eslint.Dir = o.App
	eslint.Env = append(os.Environ())
	eslint.Stderr = os.Stderr
	eslint.Stdout = os.Stdout
	eslint.Stdin = os.Stdin
	err = eslint.Run()
	if err != nil {
		return err
	}

	svelteCheck := exec.Command(bun, "x", "svelte-check", "--tsconfig=./tsconfig.json")
	svelteCheck.Dir = o.App
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
