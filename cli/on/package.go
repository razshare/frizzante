package on

import (
	"github.com/razshare/frizzante/cli"
	"github.com/razshare/frizzante/cli/path"
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"os/exec"
	"path/filepath"
)

func Package(c *cli.Cli) error {
	err := Touch(c)
	if err != nil {
		return err
	}

	bunbin, err := path.Bun(c, *c.Flags.App)
	if err != nil {
		return err
	}

	ssr := exec.Command(bunbin, "x", "vite", "build", "--logLevel=info", "--outDir=dist", "--emptyOutDir=true", "--ssr=frizzante/core/scripts/server.ts")
	ssr.Dir = *c.Flags.App
	ssr.Env = append(os.Environ())
	ssr.Stderr = os.Stderr
	ssr.Stdout = os.Stdout
	ssr.Stdin = os.Stdin
	err = ssr.Run()
	if err != nil {
		return err
	}

	csr := exec.Command(bunbin, "x", "vite", "build", "--logLevel=info", "--outDir=dist/client", "--emptyOutDir=true")
	csr.Dir = *c.Flags.App
	csr.Env = append(os.Environ())
	csr.Stderr = os.Stderr
	csr.Stdout = os.Stdout
	csr.Stdin = os.Stdin
	err = csr.Run()
	if err != nil {
		return err
	}

	esb := exec.Command("node_modules/.bin/esbuild", "--bundle", "--outfile=dist/server.js", "--format=cjs", "--allow-overwrite", "dist/server.js")
	esb.Dir = *c.Flags.App
	esb.Env = append(os.Environ())
	esb.Stderr = os.Stderr
	esb.Stdout = os.Stdout
	esb.Stdin = os.Stdin
	err = esb.Run()
	if err != nil {
		return err
	}

	messages.Success("project app package generated in ", filepath.Join(*c.Flags.App, "dist"))

	return nil
}
