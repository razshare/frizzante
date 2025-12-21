package npm

import (
	"fmt"
	"os"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
)

func Install(options InstallOptions) (err error) {
	if len(options.Packages) == 0 {
		messages.Infof("no packages to add")
		return
	}

	if !files.IsDirectory("app") {
		err = fmt.Errorf("directory %s not found", "app")
		return err
	}

	ok := 0
	for _, pkg := range options.Packages {
		messages.Infof("adding %s", pkg)
		if !messages.Command(messages.CommandOptions{
			DirectoryName: "app",
			Environment:   os.Environ(),
			Program:       options.Bun,
			Args:          []string{"add", "-D", pkg},
		}) {
			err = fmt.Errorf("failed to add package %s", pkg)
			return
		}
		messages.Successf("added %s packages to app/node_modules", pkg)
		ok++
	}

	if ok > 0 {
		messages.Successf("successfully installed %d package(s) to app/node_modules", ok)
	}

	return nil
}
