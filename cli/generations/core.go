package generations

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/confirm"
)

func Core(options CoreOptions) (err error) {
	directoryName := filepath.Join("lib", "core")
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
	directoryName = filepath.Join("app", "lib", "scripts", "core")
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
	directoryName = filepath.Join("app", "lib", "components", "core")
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

	if err = Copy(CopyOptions{
		From: "internal/project/lib/core",
		To:   filepath.Join("lib", "core"),
		Efs:  options.Efs,
	}); err != nil {
		return
	}

	if err = Copy(CopyOptions{
		From: "internal/project/app/lib/scripts/core",
		To:   filepath.Join("app", "lib", "scripts", "core"),
		Efs:  options.Efs,
	}); err != nil {
		return
	}

	if err = Copy(CopyOptions{
		From: "internal/project/app/lib/components/core",
		To:   filepath.Join("app", "lib", "components", "core"),
		Efs:  options.Efs,
	}); err != nil {
		return
	}

	if err = FixImports(FixImportsOptions{Directory: filepath.Join("lib", "core")}); err != nil {
		return
	}

	return
}
