package action

import (
	"github.com/razshare/frizzante/files"
	"os"
	"os/exec"
	"path/filepath"
)

func Check(opts CheckOptions) error {
	if err := Touch(TouchOptions{App: opts.App}); err != nil {
		return err
	}

	var bun string
	var err error
	if files.IsFile(opts.Bun) {
		if bun, err = filepath.Rel(opts.App, opts.Bun); err != nil {
			return err
		}
	} else if bun, err = exec.LookPath(opts.Bun); err != nil {
		bun = opts.Bun
	}

	eslint := exec.Command(bun, "x", "eslint")
	eslint.Dir = opts.App
	eslint.Env = append(os.Environ())
	eslint.Stderr = os.Stderr
	eslint.Stdout = os.Stdout
	eslint.Stdin = os.Stdin
	if err := eslint.Run(); err != nil {
		return err
	}

	svelteCheck := exec.Command(bun, "x", "svelte-check", "--tsconfig=./tsconfig.json")
	svelteCheck.Dir = opts.App
	svelteCheck.Env = append(os.Environ())
	svelteCheck.Stderr = os.Stderr
	svelteCheck.Stdout = os.Stdout
	svelteCheck.Stdin = os.Stdin
	if err := svelteCheck.Run(); err != nil {
		return err
	}
	return nil
}
