package codegen

import (
	"fmt"
	"github.com/razshare/frizzante/cli"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/text"
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinner"
	"os"
	"path/filepath"
)

func Download(url string) (Install, error) {
	hash, err := text.Sha1(url)
	if err != nil {
		return nil, err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	ext := filepath.Ext(url)

	if ext != ".exe" && ext != ".zip" {
		ext = ""
	}

	global := filepath.Join(home, ".frizzante", hash+ext)

	if !files.IsFile(global) {
		s := spinner.New(fmt.Sprintf("downloading %s", url))
		go spinner.Start(s)
		defer spinner.Stop(s)

		err = files.DownloadFile(url, global)
		if err != nil {
			return nil, err
		}
	}

	return func(dst string) (bool, error) {
		if files.IsDirectory(dst) || files.IsFile(dst) {
			if !*cli.Yes {
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
