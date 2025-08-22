package codegen

import (
	"fmt"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/text"
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinner"
	"os"
	"path/filepath"
)

func Download(opts DownloadOptions) (Install, Evict, error) {
	hash, err := text.Sha1(opts.Url)
	if err != nil {
		return nil, nil, err
	}

	user, err := os.UserHomeDir()
	if err != nil {
		return nil, nil, err
	}

	ext := filepath.Ext(opts.Url)

	if ext != ".exe" && ext != ".zip" {
		ext = ""
	}

	home := os.Getenv("FRIZZANTE_HOME")
	if home == "" {
		user, err = os.UserHomeDir()
		if err != nil {
			return nil, nil, err
		}
		home = filepath.Join(user, ".frizzante")
	}

	global := filepath.Join(home, hash+ext)

	if !files.IsFile(global) {
		spin := spinner.New(fmt.Sprintf("downloading %s", opts.Url))
		go spinner.Start(spin)
		defer spinner.Stop(spin)

		err = files.DownloadFile(opts.Url, global)
		if err != nil {
			return nil, nil, err
		}
	}

	return func(dst string) (bool, error) {
			if files.IsDirectory(dst) || files.IsFile(dst) {
				if !opts.Auto {
					var overwrite bool
					overwrite, err = confirm.Sendf(true, "%s already exists. Overwrite?", dst)
					if err != nil {
						return false, err
					}

					if !overwrite {
						messages.Infof("skipping %s", dst)
						return false, nil
					}
				}

				err = os.RemoveAll(dst)
				if err != nil {
					return false, err
				}
			}

			s := spinner.New(fmt.Sprintf("installing %s", dst))
			go spinner.Start(s)
			defer spinner.Stop(s)

			if ext == ".zip" {
				err = files.UnzipFile(global, dst)
				if err != nil {
					return false, err
				}
			} else {
				local := filepath.Join(dst, filepath.Base(dst)+ext)
				err = files.CopyFile(global, local)
				if err != nil {
					return false, err
				}
			}

			messages.Successf("%s installed", dst)
			return true, nil
		},
		func() error {
			return os.RemoveAll(global)
		},
		nil
}
