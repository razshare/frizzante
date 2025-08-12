package codegen

import (
	"embed"
	"github.com/razshare/frizzante/cli/state"
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"path/filepath"
)

func Links(efs embed.FS) {
	to := filepath.Join(*state.App, "frizzante", "links")

	if files.IsDirectory(to) {
		if confirm.Send(true, "feature `Links` already exists in this project. Overwrite?") {
			err := os.RemoveAll(to)
			if err != nil {
				messages.Fatal(err)
			}
		}
	}

	err := Generate(efs, []Generation{
		{
			From: "template/app/frizzante/links",
			To:   to,
		},
	})

	if err != nil {
		messages.Fatal(err)
	}
}
