package generations

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/confirm"
)

func Dev(options DevOptions) (err error) {
	directoryName := filepath.Join("lib", "dev")
	if files.IsDirectory(directoryName) {
		if options.Strict {
			err = fmt.Errorf("%s already exists", directoryName)
			return
		}
		var yesRemove bool
		if yesRemove, err = confirm.Sendf(true, "%s already exists. Remove?", directoryName); err != nil {
			return err
		}
		if !yesRemove {
			err = errors.New("cannot continue generating core library because it already exists")
			return
		}
		if err = os.RemoveAll(directoryName); err != nil {
			return
		}
	}
	if err = Copy(CopyOptions{
		From: "internal/project/lib/dev",
		To:   filepath.Join("lib", "dev"),
		Efs:  options.Efs,
	}); err != nil {
		return
	}
	if err = FixImports(FixImportsOptions{Directory: filepath.Join("lib", "dev")}); err != nil {
		return
	}
	return
}
