package action

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/razshare/frizzante/internal/tui/messages"
)

func Build(options BuildOptions) (err error) {
	if err = Package(PacakgeOptions{App: options.App, Bun: options.Bun}); err != nil {
		return
	}

	build := exec.Command(options.Go, "build", "-o="+filepath.Join(".gen", "bin", "app"), ".")
	build.Env = os.Environ()
	build.Stderr = os.Stderr
	build.Stdout = os.Stdout
	build.Stdin = os.Stdin
	if err = build.Run(); err != nil {
		return
	}

	messages.Success("project built into ", filepath.Join(".gen", "bin", "app"))

	return
}
