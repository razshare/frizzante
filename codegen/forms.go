package codegen

import (
	"embed"
	"github.com/razshare/frizzante/cli/flags"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"path/filepath"
)

func Forms(efs embed.FS) {
	to := filepath.Join(*flags.App, "frizzante", "forms")

	if files.IsDirectory(to) {
		if confirm.Send(true, "feature `Forms` already exists in this project. Overwrite?") {
			err := os.RemoveAll(to)
			if err != nil {
				messages.Fatal(err)
			}
		}
	}

	err := Generate(efs, []Generation{
		{
			From: "template/app/frizzante/forms",
			To:   to,
		},
	})

	if err != nil {
		messages.Fatal(err)
	}
}
