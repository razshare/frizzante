package actions

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/messages"
)

func Snapshot(options SnapshotOptions) (err error) {
	var name string
	if runtime.GOOS == "windows" {
		name = filepath.Join(".gen", "bin", "app.exe")
	} else {
		name = filepath.Join(".gen", "bin", "app")
	}

	if files.IsFile(filepath.Join(".gen", "bin", "app")) {
		if !options.Auto {
			var yes bool
			if yes, err = confirm.Sendf(false, "file %s already exists. Rebuild?", name); err != nil {
				return
			}

			if yes {
				if err = os.RemoveAll(filepath.Join(".gen", "bin", "app")); err != nil {
					return
				}

				if err = Build(BuildOptions{
					Go:       options.Go,
					Bun:      options.Bun,
					Tags:     options.Tags,
					Platform: options.Platform,
				}); err != nil {
					return
				}
			} else {
				messages.Info("skipping build")
			}
		}
	}

	return
}
