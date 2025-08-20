package codegen

import (
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/messages"
	"os"
)

func Forms(o FormsOptions) error {
	if files.IsDirectory(o.Lib) {
		if !o.Auto {
			overwrite, err := confirm.Sendf(true, "%s already exists. Overwrite?", o.Lib)
			if err != nil {
				return err
			}

			if !overwrite {
				messages.Infof("skipping %s", o.Lib)
				return nil
			}
		}

		err := os.RemoveAll(o.Lib)
		if err != nil {
			return err
		}
	}

	err := Copy([]CopyInstruction{
		{
			Efs:  o.Efs,
			From: "template/app/frizzante/forms",
			To:   o.Lib,
		},
	})

	if err != nil {
		return err
	}

	messages.Successf("form files generated at %s", o.Lib)

	return nil
}
