package on

import (
	"github.com/razshare/frizzante/cli/flags"
	"github.com/razshare/frizzante/cli/path"
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"os/exec"
)

func Check() {
	Touch()

	eslint := exec.Command(path.Bun(*flags.App), "x", "eslint")
	eslint.Dir = *flags.App
	eslint.Env = append(os.Environ())
	eslint.Stderr = os.Stderr
	eslint.Stdout = os.Stdout
	eslint.Stdin = os.Stdin
	err := eslint.Run()
	if err != nil {
		messages.Fatal(err)
	}

	svelteCheck := exec.Command(path.Bun(*flags.App), "x", "svelte-check", "--tsconfig=./tsconfig.json")
	svelteCheck.Dir = *flags.App
	svelteCheck.Env = append(os.Environ())
	svelteCheck.Stderr = os.Stderr
	svelteCheck.Stdout = os.Stdout
	svelteCheck.Stdin = os.Stdin
	err = svelteCheck.Run()
	if err != nil {
		messages.Fatal(err)
	}
}
