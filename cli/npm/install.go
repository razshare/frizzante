package npm

import (
	"fmt"
	"os"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
)

func Install(bun string, app string, pkgs ...string) error {
	if len(pkgs) == 0 {
		return nil
	}

	if !files.IsDirectory(app) {
		return fmt.Errorf("directory %s not found", app)
	}

	ok := 0
	for _, pkg := range pkgs {
		messages.Infof("adding %s", pkg)
		if !messages.Command(app, os.Environ(), bun, "add", "-D", pkg) {
			messages.Errorf("failed to add package %s", pkg)
			continue
		}
		messages.Successf("added %s packages to %s/node_modules", pkg, app)
		ok++
	}

	if ok > 0 {
		messages.Successf("successfully installed %d package(s) to %s/node_modules", ok, app)
	}

	return nil
}
