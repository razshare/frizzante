package codegen

import (
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"path/filepath"
)

func Core(app string, yes bool) error {
	dst := filepath.Join(app, "frizzante", "core")

	if files.IsDirectory(dst) {
		if !yes {
			overwrite, err := confirm.Sendf(true, "%s already exists. Overwrite?", dst)
			if err != nil {
				return err
			}

			if !overwrite {
				messages.Infof("skipping %s", dst)
				return nil
			}
		}

		err := os.RemoveAll(dst)
		if err != nil {
			return err
		}
	}

	err := Copy([]CopyInstruction{
		{
			From: "template/app/frizzante/core",
			To:   dst,
		},
	})

	if err != nil {
		return err
	}

	messages.Successf("core files generated at %s", dst)

	return nil
}
