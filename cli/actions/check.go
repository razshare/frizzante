package actions

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinners"
)

func Check(options CheckOptions) (err error) {
	spin := spinners.New("checking code")
	go spinners.Start(spin)

	if err = Touch(TouchOptions{}); err != nil {
		spinners.Stop(spin)
		return
	}

	var bun string
	if files.IsFile(options.Bun) {
		if bun, err = filepath.Rel("app", options.Bun); err != nil {
			spinners.Stop(spin)
			return
		}
	} else if bun, err = exec.LookPath(options.Bun); err != nil {
		spinners.Stop(spin)
		bun = options.Bun
	}

	eslint := exec.Command(bun, "x", "eslint")
	eslint.Dir = "app"
	eslint.Env = os.Environ()
	eslint.Stderr = os.Stderr
	eslint.Stdout = os.Stdout
	eslint.Stdin = os.Stdin
	if err = eslint.Run(); err != nil {
		spinners.Stop(spin)
		if eslint.Err != nil {
			messages.Error(eslint.Err.Error())
		}
		return
	}

	var data []byte
	if data, err = os.ReadFile(filepath.Join("app", "package.json")); err != nil {
		spinners.Stop(spin)
		return
	}

	type DevDependencies struct {
		SvelteCheck string `json:"svelte-check"`
	}

	type PackageJson struct {
		DevDependencies DevDependencies `json:"devDependencies"`
	}

	var pkg PackageJson
	if err = json.Unmarshal(data, &pkg); err != nil {
		spinners.Stop(spin)
		return
	}

	if pkg.DevDependencies.SvelteCheck != "" {
		svelteCheck := exec.Command(bun, "x", "svelte-check", "--tsconfig=./tsconfig.json")
		svelteCheck.Dir = "app"
		svelteCheck.Env = os.Environ()
		svelteCheck.Stderr = os.Stderr
		svelteCheck.Stdout = os.Stdout
		svelteCheck.Stdin = os.Stdin
		err = svelteCheck.Run()
		if svelteCheck.Err != nil {
			spinners.Stop(spin)
			messages.Error(eslint.Err.Error())
			return
		}
	}

	spinners.Stop(spin)

	return
}
