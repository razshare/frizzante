package actions

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinners"
)

func Build(options BuildOptions) (err error) {
	if err = Package(PackageOptions{Bun: options.Bun, Prod: true}); err != nil {
		return
	}

	if len(options.Tags) > 0 {
		spin := spinners.Newf("building binary with tags %s", strings.Join(options.Tags, ","))
		go spinners.Start(spin)
		defer spinners.Stop(spin)
		if messages.Command("", os.Environ(), options.Go, "build", "-tags="+strings.Join(options.Tags, ","), "-o="+filepath.Join(".gen", "bin", "app"), ".") {
			messages.Success("project built into ", filepath.Join(".gen", "bin", "app"))
		}
	} else {
		spin := spinners.New("building binary")
		go spinners.Start(spin)
		defer spinners.Stop(spin)
		if messages.Command("", os.Environ(), options.Go, "build", "-o="+filepath.Join(".gen", "bin", "app"), ".") {
			messages.Success("project built into ", filepath.Join(".gen", "bin", "app"))
		}
	}

	return
}
