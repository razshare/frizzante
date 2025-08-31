package generate

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/text"
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinner"
)

func Download(options DownloadOptions) (ins Install, evc Evict, err error) {
	var hash string
	if hash, err = text.Sha1(options.Url); err != nil {
		return
	}

	var user string
	if user, err = os.UserHomeDir(); err != nil {
		return
	}

	ext := filepath.Ext(options.Url)

	if ext != ".exe" && ext != ".zip" {
		ext = ""
	}

	home := os.Getenv("FRIZZANTE_HOME")
	if home == "" {
		if user, err = os.UserHomeDir(); err != nil {
			return
		}
		home = filepath.Join(user, ".frizzante")
	}

	global := filepath.Join(home, hash+ext)

	if !files.IsFile(global) {
		spin := spinner.New(fmt.Sprintf("downloading %s", options.Url))
		go spinner.Start(spin)
		defer spinner.Stop(spin)
		if err = files.DownloadFile(options.Url, global); err != nil {
			return nil, nil, err
		}
	}

	return func(dst string) (installed bool, err error) {
			if files.IsDirectory(dst) || files.IsFile(dst) {
				if !options.Auto {
					var overwrite bool
					if overwrite, err = confirm.Sendf(true, "%s already exists. Overwrite?", dst); err != nil {
						return
					}

					if !overwrite {
						messages.Infof("skipping %s", dst)
						return
					}
				}

				if err = os.RemoveAll(dst); err != nil {
					return
				}
			}

			spin := spinner.New(fmt.Sprintf("installing %s", dst))
			go spinner.Start(spin)
			defer spinner.Stop(spin)

			if ext == ".zip" {
				if err = files.UnzipFile(global, dst); err != nil {
					return
				}
			} else {
				local := filepath.Join(dst, filepath.Base(dst)+ext)
				if err = files.CopyFile(global, local); err != nil {
					return
				}
			}

			installed = true

			messages.Successf("%s installed", dst)

			return
		},
		func() error { return os.RemoveAll(global) },
		nil
}
