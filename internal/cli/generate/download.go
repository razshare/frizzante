package generate

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/internal/cli/user"
	"github.com/razshare/frizzante/internal/text"
	"github.com/razshare/frizzante/tui/confirm"
	messages2 "github.com/razshare/frizzante/tui/messages"
	spinner2 "github.com/razshare/frizzante/tui/spinner"
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

	if !files.IsFile(global) {
		spin := spinner2.New(fmt.Sprintf("downloading %s", options.Url))
		go spinner2.Start(spin)
		defer spinner2.Stop(spin)
		if err = files.DownloadFile(options.Url, global); err != nil {
			return nil, nil, err
		}
	}

	return func(to string) (installed bool, err error) {
			if files.IsDirectory(to) || files.IsFile(to) {
				if !options.Auto {
					var overwrite bool
					if overwrite, err = confirm.Sendf(true, "%s already exists. Overwrite?", to); err != nil {
						return
					}

					if !overwrite {
						messages2.Infof("skipping %s", to)
						return
					}
				}

				if err = os.RemoveAll(to); err != nil {
					return
				}
			}

			spin := spinner2.New(fmt.Sprintf("installing %s", to))
			go spinner2.Start(spin)
			defer spinner2.Stop(spin)

			if ext == ".zip" {
				if err = files.UnzipFile(global, to); err != nil {
					return
				}
			} else {
				local := filepath.Join(to, filepath.Base(to)+ext)
				if err = files.CopyFile(global, local); err != nil {
					return
				}
			}

			installed = true

			messages2.Successf("%s installed", to)

			return
		},
		func() error { return os.RemoveAll(global) },
		nil
}
