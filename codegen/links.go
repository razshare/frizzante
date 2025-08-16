package codegen

import (
	"github.com/razshare/frizzante/cli"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"path/filepath"
)

func Links(c *cli.Cli, _ bool, base string) error {
	dst := filepath.Join(base, *c.Flags.App, "frizzante", "links")

	if files.IsDirectory(dst) {
		overwrite, err := confirm.Sendf(true, "%s already exists. Overwrite?", dst)
		if err != nil {
			return err
		}

		if !overwrite {
			messages.Infof("skipping %s", dst)
			return nil
		}

		err = os.RemoveAll(dst)
		if err != nil {
			return err
		}
	}

	err := Copy(c.Efs, []CopyInstruction{
		{
			From: "template/app/frizzante/links",
			To:   dst,
		},
	})

	if err != nil {
		return err
	}

	messages.Successf("link files generated at %s", dst)

	return nil
}
