package generate

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/razshare/frizzante/cli/paths"
	"github.com/razshare/frizzante/internal/additions/lib/security"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinners"
)

func Download(options DownloadOptions) (install Install, evict Evict, err error) {
	var cache string
	if cache, err = paths.Cache(); err != nil {
		return
	}

	hash := security.Sha1(options.Url)

	ext := filepath.Ext(options.Url)

	if ext != ".exe" && ext != ".zip" {
		ext = ""
	}

	global := filepath.Join(cache, hash+ext)

	if !files.IsFile(global) {
		spin := spinners.New(fmt.Sprintf("downloading %s", options.Url))
		go spinners.Start(spin)
		if err = files.DownloadFile(options.Url, global); err != nil {
			spinners.Stop(spin)
			return
		}
		spinners.Stop(spin)
	}

	install = func(to string) (installed bool, err error) {
		if files.IsDirectory(to) || files.IsFile(to) {
			if !options.Auto {
				var overwrite bool
				if overwrite, err = confirm.Sendf(true, "%s already exists. Overwrite?", to); err != nil {
					return
				}

				if !overwrite {
					messages.Infof("skipping %s", to)
					return
				}
			}

			if err = os.RemoveAll(to); err != nil {
				return
			}
		}

		spin := spinners.New(fmt.Sprintf("installing %s", to))
		go spinners.Start(spin)

		if ext == ".zip" {
			if err = files.UnzipFile(global, to); err != nil {
				spinners.Stop(spin)
				return
			}
		} else {
			local := filepath.Join(to, filepath.Base(to)+ext)
			if err = files.CopyFile(global, local); err != nil {
				spinners.Stop(spin)
				return
			}
		}

		spinners.Stop(spin)

		installed = true

		messages.Successf("%s installed", to)

		return
	}

	evict = func() error { return os.RemoveAll(global) }

	return
}
