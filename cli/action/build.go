package action

import (
	"os"
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

	if messages.Command(".", os.Environ(), options.Go, "build", "-o="+filepath.Join(".gen", "bin", "app"), ".") {
		messages.Success("project built into ", filepath.Join(".gen", "bin", "app"))
	}

	return
}
