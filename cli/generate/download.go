package generate

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/razshare/frizzante/cli/user"
	files2 "github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/text"
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinner"
)

func Download(options DownloadOptions) (install Install, evict Evict, err error) {
	var cache string
	if cache, err = user.FrizzanteCache(); err != nil {
		return
	}

	var hash string
	if hash, err = text.Sha1(options.Url); err != nil {
		return
	}

	ext := filepath.Ext(options.Url)

	if ext != ".exe" && ext != ".zip" {
		ext = ""
	}

	global := filepath.Join(cache, hash+ext)

	if !files2.IsFile(global) {
		spin := spinner.New(fmt.Sprintf("downloading %s", options.Url))
		go spinner.Start(spin)
		defer spinner.Stop(spin)
		if err = files2.DownloadFile(options.Url, global); err != nil {
			return nil, nil, err
		}
	}

	return func(to string) (installed bool, err error) {
			if files2.IsDirectory(to) || files2.IsFile(to) {
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

			spin := spinner.New(fmt.Sprintf("installing %s", to))
			go spinner.Start(spin)
			defer spinner.Stop(spin)

			if ext == ".zip" {
				if err = files2.UnzipFile(global, to); err != nil {
					return
				}
			} else {
				local := filepath.Join(to, filepath.Base(to)+ext)
				if err = files2.CopyFile(global, local); err != nil {
					return
				}
			}

			installed = true

			messages.Successf("%s installed", to)

			return
		},
		func() error { return os.RemoveAll(global) },
		nil
}
