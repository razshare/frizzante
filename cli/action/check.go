package action

import (
	"github.com/razshare/frizzante/files"
	"os"
	"os/exec"
	"path/filepath"
)

func Check(options CheckOptions) (err error) {
	if err = Touch(TouchOptions{App: options.App}); err != nil {
		return
	}

	var bun string
	if files.IsFile(options.Bun) {
		if bun, err = filepath.Rel(options.App, options.Bun); err != nil {
			return
		}
	} else if bun, err = exec.LookPath(options.Bun); err != nil {
		bun = options.Bun
	}

	eslint := exec.Command(bun, "x", "eslint")
	eslint.Dir = options.App
	eslint.Env = append(os.Environ())
	eslint.Stderr = os.Stderr
	eslint.Stdout = os.Stdout
	eslint.Stdin = os.Stdin
	if err = eslint.Run(); err != nil {
		return
	}

	svelteCheck := exec.Command(bun, "x", "svelte-check", "--tsconfig=./tsconfig.json")
	svelteCheck.Dir = options.App
	svelteCheck.Env = append(os.Environ())
	svelteCheck.Stderr = os.Stderr
	svelteCheck.Stdout = os.Stdout
	svelteCheck.Stdin = os.Stdin

	return svelteCheck.Run()
}
