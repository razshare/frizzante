package codegen

import (
	"embed"
	"github.com/razshare/frizzante/cli/state"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/tui/confirm"
	"os"
	"path/filepath"
)

func Core(efs embed.FS, base string) error {
	dst := filepath.Join(base, *state.App, "frizzante", "core")

	if files.IsDirectory(dst) {
		overwrite, err := confirm.Sendf(true, "%s already exists. Overwrite?", dst)
		if err != nil {
			return err
		}

		if overwrite {
			err = os.RemoveAll(dst)
			if err != nil {
				return err
			}
		}
	}

	err := Copy(efs, []Generation{
		{
			From: "template/app/frizzante/core",
			To:   dst,
		},
	})

	if err != nil {
		return err
	}

	return nil
}
