package generations

import (
	"errors"
	"fmt"
	"os"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/confirm"
)

func Makefile(options MakefileOptions) (err error) {
	fileName := "makefile"
	if files.IsFile(fileName) {
		if options.Strict {
			err = fmt.Errorf("%s already exists", fileName)
			return
		}
		var yesRemove bool
		if yesRemove, err = confirm.Sendf(true, "%s already exists. Remove?", fileName); err != nil {
			return err
		}
		if !yesRemove {
			err = errors.New("cannot continue generating core library because it already exists")
			return
		}
		if err = os.RemoveAll(fileName); err != nil {
			return
		}
	}
	if err = Copy(CopyOptions{
		From: fmt.Sprintf("internal/project/%s", fileName),
		To:   fileName,
		Efs:  options.Efs,
	}); err != nil {
		return
	}
	return
}
