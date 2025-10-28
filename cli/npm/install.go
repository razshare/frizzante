package npm

import (
	"fmt"
	"os"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
)

func Install(options InstallOptions) error {
	if len(options.Packages) == 0 {
		return nil
	}

	if !files.IsDirectory("app") {
		return fmt.Errorf("directory %s not found", "app")
	}

	ok := 0
	for _, pkg := range options.Packages {
		messages.Infof("adding %s", pkg)
		if !messages.Command(messages.CommandOptions{
			Dir:  "app",
			Env:  os.Environ(),
			Name: options.Bun,
			Args: []string{"add", "-D", pkg},
		}) {
			messages.Errorf("failed to add package %s", pkg)
			continue
		}
		messages.Successf("added %s packages to app/node_modules", pkg)
		ok++
	}

	if ok > 0 {
		messages.Successf("successfully installed %d package(s) to app/node_modules", ok)
	}

	return nil
}
