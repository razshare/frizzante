package codegen

import (
	"embed"
	"github.com/razshare/frizzante/cli/state"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/tui/confirm"
	"os"
	"path/filepath"
)

func Links(efs embed.FS, base string) error {
	dst := filepath.Join(base, *state.App, "frizzante", "links")

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

	return Copy(efs, []Generation{
		{
			From: "template/app/frizzante/links",
			To:   dst,
		},
	})
}
