package action

import (
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"os/exec"
	"path/filepath"
)

func Build(o BuildOptions) error {
	err := Pkg(PkgOptions{App: o.App, Bun: o.Bun})
	if err != nil {
		return err
	}

	build := exec.Command(o.Go, "build", "-o="+filepath.Join(".gen", "bin", "app"), ".")
	build.Env = os.Environ()

	build.Stderr = os.Stderr
	build.Stdout = os.Stdout
	build.Stdin = os.Stdin
	err = build.Run()
	if err != nil {
		return err
	}

	messages.Success("project built into ", filepath.Join(".gen", "bin", "app"))

	return nil
}
