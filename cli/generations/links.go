package generations

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/confirm"
)

func Links(options LinksOptions) (err error) {
	directoryName := filepath.Join("app", "lib", "components", "links")
	if files.IsDirectory(directoryName) {
		if options.Strict {
			err = fmt.Errorf("%s already exists", directoryName)
			return
		}
		var yesRemove bool
		if yesRemove, err = confirm.Sendf(true, "%s already exists. Remove?", directoryName); err != nil {
			return err
		}
		if yesRemove {
			if err = os.RemoveAll(directoryName); err != nil {
				return
			}
		}
	}
	err = Copy(CopyOptions{
		From: "internal/additions/app/lib/components/links",
		To:   filepath.Join("app", "lib", "components", "links"),
		Efs:  options.Efs,
	})
	return
}
