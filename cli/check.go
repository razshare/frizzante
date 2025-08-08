package cli

import (
	"os"
	"os/exec"
)

func OnCheck() {
	OnTouch()

	eslint := exec.Command(Bun("app"), "x", "eslint")
	eslint.Dir = "app"
	eslint.Env = append(os.Environ())
	eslint.Stderr = os.Stderr
	eslint.Stdout = os.Stdout
	eslint.Stdin = os.Stdin
	eslintError := eslint.Run()
	if eslintError != nil {
		Fatal(eslintError)
	}

	svelteCheck := exec.Command(Bun("app"), "x", "svelte-check", "--tsconfig=./tsconfig.json")
	svelteCheck.Dir = "app"
	svelteCheck.Env = append(os.Environ())
	svelteCheck.Stderr = os.Stderr
	svelteCheck.Stdout = os.Stdout
	svelteCheck.Stdin = os.Stdin
	svelteCheckError := svelteCheck.Run()
	if svelteCheckError != nil {
		Fatal(svelteCheckError)
	}
}
