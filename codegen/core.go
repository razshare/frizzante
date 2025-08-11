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

func Core(efs embed.FS) {
	to := filepath.Join(*flags.App, "frizzante", "core")

	if files.IsDirectory(to) {
		if confirm.Send(true, "feature `Core` already exists in this project. Overwrite?") {
			err := os.RemoveAll(to)
			if err != nil {
				messages.Fatal(err)
			}
		}
	}

	err := Generate(efs, []Generation{
		{
			From: "template/app/frizzante/core",
			To:   to,
		},
	})

	if err != nil {
		messages.Fatal(err)
	}
}
