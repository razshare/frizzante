package action

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/tui/messages"
)

func Test(options TestOptions) (err error) {
	if files.IsDirectory(filepath.Join("lib", "core", "view", "ssr")) {
		if files.IsDirectory(filepath.Join(options.App, "dist")) {
			if err = files.CopyDirectory(
				filepath.Join(options.App, "dist"),
				filepath.Join("lib", "core", "view", "ssr", "app", "dist"),
			); err != nil {
				return
			}

			messages.Successf("%s copied to  %s", filepath.Join(options.App, "dist"), filepath.Join("lib", "core", "view", "ssr"))
		}
	}

	test := exec.Command(options.Go, "test", "./...")
	test.Env = os.Environ()
	test.Stderr = os.Stderr
	test.Stdout = os.Stdout
	test.Stdin = os.Stdin
	err = test.Run()

	return
}
