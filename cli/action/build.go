package action

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinner"
)

func Build(options BuildOptions) (err error) {
	var spin *spinner.Spinner

	if len(options.Tags) > 0 {
		spin = spinner.Newf("building binary with tags %s", strings.Join(options.Tags, ","))
	} else {
		spin = spinner.New("building binary")
	}

	go spinner.Start(spin)
	defer spinner.Stop(spin)

	if err = Package(PackageOptions{App: options.App, Bun: options.Bun, Prod: true}); err != nil {
		return
	}

	if len(options.Tags) > 0 {
		if messages.Command(".", os.Environ(), options.Go, "build", "-tags="+strings.Join(options.Tags, ","), "-o="+filepath.Join(".gen", "bin", "app"), ".") {
			messages.Success("project built into ", filepath.Join(".gen", "bin", "app"))
		}
	} else {
		if messages.Command(".", os.Environ(), options.Go, "build", "-o="+filepath.Join(".gen", "bin", "app"), ".") {
			messages.Success("project built into ", filepath.Join(".gen", "bin", "app"))
		}
	}

	return
}
