package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/razshare/frizzante/embeds"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/messages"
)

func Copy(options CopyOptions) (err error) {
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

	if embeds.IsDirectory(options.Efs, options.From) {
		var entries []string
		if entries, err = embeds.ReadDirectory(options.Efs, options.From); err != nil {
			return err
		}

		for _, entry := range entries {
			if options.Ignore != nil {
				var ignored bool
				for _, ignore := range options.Ignore {
					ignored = strings.HasPrefix(entry, ignore)
					if ignored {
						break
					}
				}

				if ignored {
					continue
				}
			}

			name := filepath.Join(options.To, strings.TrimPrefix(entry, options.From))
			err = embeds.CopyFile(options.Efs, entry, name)
			if err != nil {
				return err
			}
		}
	} else if embeds.IsFile(options.Efs, options.From) {
		if err = embeds.CopyFile(options.Efs, options.From, options.To); err != nil {
			return
		}
	} else {
		err = fmt.Errorf("%s not found", options.From)
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
