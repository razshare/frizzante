package action

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinner"
)

func Check(options CheckOptions) (err error) {
	spin := spinner.New("checking code")
	go spinner.Start(spin)
	defer spinner.Stop(spin)

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
	eslint.Env = os.Environ()
	//eslint.Stderr = os.Stderr
	//eslint.Stdout = os.Stdout
	//eslint.Stdin = os.Stdin
	if err = eslint.Run(); err != nil {
		if eslint.Err != nil {
			messages.Error(eslint.Err.Error())
		}
		return
	}

	var data []byte
	if data, err = os.ReadFile(filepath.Join(options.App, "package.json")); err != nil {
		return err
	}

	type DevDependencies struct {
		SvelteCheck string `json:"svelte-check"`
	}

	type PackageJson struct {
		DevDependencies DevDependencies `json:"devDependencies"`
	}

	var pkg PackageJson
	if err = json.Unmarshal(data, &pkg); err != nil {
		return
	}

	if pkg.DevDependencies.SvelteCheck != "" {
		svelteCheck := exec.Command(bun, "x", "svelte-check", "--tsconfig=./tsconfig.json")
		svelteCheck.Dir = options.App
		svelteCheck.Env = os.Environ()
		//svelteCheck.Stderr = os.Stderr
		//svelteCheck.Stdout = os.Stdout
		//svelteCheck.Stdin = os.Stdin
		err = svelteCheck.Run()
		if svelteCheck.Err != nil {
			messages.Error(eslint.Err.Error())
		}
	}

	return
}
