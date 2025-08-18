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

func Download(o DownloadOptions) (Install, error) {
	hash, err := text.Sha1(o.Url)
	if err != nil {
		return nil, err
	}

	user, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	ext := filepath.Ext(o.Url)

	if ext != ".exe" && ext != ".zip" {
		ext = ""
	}

	home := os.Getenv("FRIZZANTE_HOME")
	if home == "" {
		user, err = os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		home = filepath.Join(user, ".frizzante")
	}

	global := filepath.Join(home, hash+ext)

	if !files.IsFile(global) {
		s := spinner.New(fmt.Sprintf("downloading %s", o.Url))
		go spinner.Start(s)
		defer spinner.Stop(s)

		err = files.DownloadFile(o.Url, global)
		if err != nil {
			return nil, err
		}
	}

	return func(dst string) (bool, error) {
		if files.IsDirectory(dst) || files.IsFile(dst) {
			if !o.Auto {
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
	}, nil
}
