package action

import (
	"os"
	"os/exec"
	"path/filepath"

	files2 "github.com/razshare/frizzante/files"
)

func Test(options TestOptions) (err error) {
	if files2.IsFile(filepath.Join("lib", "core", "svelte", "ssr")) && !files2.IsDirectory(filepath.Join("lib", "core", "svelte", "ssr", "app")) {
		if !files2.IsDirectory(filepath.Join(options.App, "dist")) && files2.IsDirectory(filepath.Join(options.App, "dist")) {
			if err = files2.CopyDirectory(
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
