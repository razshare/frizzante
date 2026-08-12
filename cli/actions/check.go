package actions

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/razshare/frizzante/v2/tui/messages"
	"github.com/razshare/frizzante/v2/tui/spinners"
)

func Check(options CheckOptions) (err error) {
	spin := spinners.New("checking code")
	go spinners.Start(spin)
	defer spinners.Stop(spin)
	var stdErrBuilder strings.Builder
	if !messages.Command(messages.CommandOptions{
		Environment:   os.Environ(),
		Program:       options.Go,
		StderrBuilder: &stdErrBuilder,
		DisableStderr: true,
		Args:          []string{"vet"},
	}) {
		err = fmt.Errorf("go vet failed:%s", stdErrBuilder.String())
		return
	}
	if !messages.Command(messages.CommandOptions{
		DirectoryName: "app",
		Environment:   os.Environ(),
		Program:       options.Bun,
		StderrBuilder: &stdErrBuilder,
		DisableStderr: true,
		Args:          []string{"x", "eslint"},
	}) {
		err = fmt.Errorf("eslint failed:%s", stdErrBuilder.String())
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
	args := []string{"x", "svelte-check", "--tsconfig=./tsconfig.json"}
	if options.Incremental {
		args = append(args, "--incremental")
	}
	if pkg.DevDependencies.SvelteCheck != "" {
		if !messages.Command(messages.CommandOptions{
			DirectoryName: "app",
			Environment:   os.Environ(),
			Program:       options.Bun,
			StderrBuilder: &stdErrBuilder,
			DisableStderr: true,
			Args:          args,
		}) {
			err = fmt.Errorf("svelte-check failed:%s", stdErrBuilder.String())
			return
		}
	}
	return
}
