package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/internal/cli/user"
	"github.com/razshare/frizzante/internal/embeds"
	"github.com/razshare/frizzante/internal/tui/confirm"
	messages2 "github.com/razshare/frizzante/internal/tui/messages"
)

func Copy(options CopyOptions) (err error) {
	var home string
	if home, err = user.FrizzanteHome(); err != nil {
		return
	}

	var data []byte
	if data, err = options.Efs.ReadFile("version"); err != nil {
		return
	}

	version := string(data)

	if !files.IsFile(filepath.Join(home, "project-"+version+".zip")) || !files.IsDirectory(filepath.Join(home, "project")) {
		if err = os.RemoveAll(filepath.Join(home, "project")); err != nil {
			return
		}

		if err = os.RemoveAll(filepath.Join(home, "project-"+version+".zip")); err != nil {
			return
		}

		if err = embeds.CopyFile(options.Efs, "internal/project.zip", filepath.Join(home, "project-"+version+".zip")); err != nil {
			return
		}

		if err = files.UnzipFile(filepath.Join(home, "project-"+version+".zip"), home); err != nil {
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
				messages2.Infof("skipping %s", options.To)
				return
			}
		}

		if err = os.RemoveAll(options.To); err != nil {
			return
		}
	}

	if files.IsDirectory(filepath.Join(home, "project", options.From)) {
		var entries []string
		if entries, err = files.ReadDirectory(filepath.Join(home, "project", options.From)); err != nil {
			return
		}

		for _, entry := range entries {
			entryRelative := strings.TrimPrefix(entry, home+string(filepath.Separator)+"project"+string(filepath.Separator))
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
	} else if files.IsFile(filepath.Join(home, "project", options.From)) {
		if err = files.CopyFile(filepath.Join(home, "project", options.From), options.To); err != nil {
			return
		}
	} else {
		err = fmt.Errorf("%s not found", filepath.Join(home, "project", options.From))
		return
	}

	if options.Ignore != nil {
		for _, name := range options.Ignore {
			if err = os.RemoveAll(name); err != nil {
				return
			}
		}
	}

	messages2.Successf("%s created", options.To)

	return
}
