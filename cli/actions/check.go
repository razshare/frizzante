package actions

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinners"
)

func Check(options CheckOptions) (err error) {
	spin := spinners.New("checking code")
	go spinners.Start(spin)
	defer spinners.Stop(spin)
	if !messages.Command(messages.CommandOptions{
		DirectoryName: "app",
		Environment:   os.Environ(),
		Program:       options.Bun,
		Args:          []string{"x", "eslint"},
	}) {
		err = errors.New("could not run eslint")
		return
	}
	var data []byte
	if data, err = os.ReadFile(filepath.Join("app", "package.json")); err != nil {
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
		return
	}
	if pkg.DevDependencies.SvelteCheck != "" {
		if !messages.Command(messages.CommandOptions{
			DirectoryName: "app",
			Environment:   os.Environ(),
			Program:       options.Bun,
			Args:          []string{"x", "svelte-check", "--tsconfig=./tsconfig.json"},
		}) {
			err = errors.New("could not run svelte-check")
			return
		}
	}
	return
}
