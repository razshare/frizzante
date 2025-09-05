package action

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/razshare/frizzante/internal/files"
)

func Test(options TestOptions) (err error) {
	if files.IsFile(filepath.Join("lib", "core", "svelte", "ssr")) && !files.IsDirectory(filepath.Join("lib", "core", "svelte", "ssr", "app")) {
		if !files.IsDirectory(filepath.Join(options.App, "dist")) && files.IsDirectory(filepath.Join(options.App, "dist")) {
			if err = files.CopyDirectory(
				filepath.Join(options.App, "dist"),
				filepath.Join("lib", "core", "svelte", "ssr", "app", "dist"),
			); err != nil {
				return
			}
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
