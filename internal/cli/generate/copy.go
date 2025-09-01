package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/razshare/frizzante/embeds"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/internal/cli/user"
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/messages"
)

func Copy(options CopyOptions) (err error) {
	var cache string
	if cache, err = user.FrizzanteCache(); err != nil {
		return
	}

	var data []byte
	if data, err = options.Efs.ReadFile("version"); err != nil {
		return
	}

	version := string(data)

	if !files.IsFile(filepath.Join(cache, "project-"+version+".zip")) || !files.IsDirectory(filepath.Join(cache, "project")) {
		if err = os.RemoveAll(filepath.Join(cache, "project")); err != nil {
			return
		}

		if err = os.RemoveAll(filepath.Join(cache, "project-"+version+".zip")); err != nil {
			return
		}

		if err = embeds.CopyFile(options.Efs, "internal/project.zip", filepath.Join(cache, "project-"+version+".zip")); err != nil {
			return
		}

		if err = files.UnzipFile(filepath.Join(cache, "project-"+version+".zip"), cache); err != nil {
			return
		}
	}

	if files.IsFile(options.To) || files.IsDirectory(options.To) {
		if !options.Auto {
			var overwrite bool
			if overwrite, err = confirm.Sendf(true, "%s already exists. Overwrite?", options.To); err != nil {
				return
			}

			if !overwrite {
				messages.Infof("skipping %s", options.To)
				return
			}
		}

		if err = os.RemoveAll(options.To); err != nil {
			return
		}
	}

	if files.IsDirectory(filepath.Join(cache, "project", options.From)) {
		var entries []string
		if entries, err = files.ReadDirectory(filepath.Join(cache, "project", options.From)); err != nil {
			return
		}

		for _, entry := range entries {
			entryRelative := strings.TrimPrefix(entry, cache+string(filepath.Separator)+"project"+string(filepath.Separator))
			if options.Ignore != nil {
				var ignored bool
				for _, ignore := range options.Ignore {
					ignored = strings.HasPrefix(entryRelative, ignore)
					if ignored {
						break
					}
				}

				if ignored {
					continue
				}
			}

			name := filepath.Join(options.To, strings.TrimPrefix(entryRelative, options.From))

			if err = files.CopyFile(entry, name); err != nil {
				return
			}
		}
	} else if files.IsFile(filepath.Join(cache, "project", options.From)) {
		if err = files.CopyFile(filepath.Join(cache, "project", options.From), options.To); err != nil {
			return
		}
	} else {
		err = fmt.Errorf("%s not found", filepath.Join(cache, "project", options.From))
		return
	}

	if options.Ignore != nil {
		for _, name := range options.Ignore {
			if err = os.RemoveAll(name); err != nil {
				return
			}
		}
	}

	messages.Successf("%s created", options.To)

	return
}
