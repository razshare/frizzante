package action

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinner"
)

func Build(options BuildOptions) (err error) {
	spin := spinner.New("building binary")
	go spinner.Start(spin)
	defer spinner.Stop(spin)

	if err = Package(PackageOptions{App: options.App, Bun: options.Bun, Prod: true}); err != nil {
		return
	}

	build := exec.Command(options.Go, "build", "-o="+filepath.Join(".gen", "bin", "app"), ".")
	build.Env = os.Environ()
	//build.Stderr = os.Stderr
	//build.Stdout = os.Stdout
	//build.Stdin = os.Stdin
	if err = build.Run(); err != nil {
		if build.Err != nil {
			messages.Error(build.Err.Error())
		}
		return
	}

	messages.Success("project built into ", filepath.Join(".gen", "bin", "app"))

	return
}
