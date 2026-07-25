package npm

import (
	"fmt"
	"os"

	"github.com/razshare/frizzante/v2/internal/project/lib/core/files"
	"github.com/razshare/frizzante/v2/tui/messages"
)

func Install(options InstallOptions) (err error) {
	if len(options.PackageNames) == 0 {
		messages.Infof("no packages to add")
		return
	}
	if !files.IsDirectory("app") {
		err = fmt.Errorf("directory %s not found", "app")
		return err
	}
	ok := 0
	for _, packageName := range options.PackageNames {
		messages.Infof("adding %s", packageName)
		if !messages.Command(messages.CommandOptions{
			DirectoryName: "app",
			Environment:   os.Environ(),
			Program:       options.Bun,
			Args:          []string{"add", "-D", packageName},
		}) {
			err = fmt.Errorf("failed to add package %s", packageName)
			return
		}
		messages.Successf("added %s packages to app/node_modules", packageName)
		ok++
	}
	if ok > 0 {
		messages.Successf("successfully installed %d package(s) to app/node_modules", ok)
	}
	return
}
