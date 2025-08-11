package copy_modded

import (
	"embed"
	"github.com/razshare/frizzante/codegen"
	"github.com/razshare/frizzante/embeds"
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/messages"
	"path/filepath"
)

func Database(efs embed.FS) {
	err := embeds.Generate(efs, []codegen.Generation{
		{
			From: "template/lib/database",
			To:   filepath.Join("lib", "database"),
			Overwrite: func(n string) bool {
				yes := confirm.Sendf(true, "file `%s` already exists. Overwrite?", n)
				if yes {
					messages.Infof("overwriting file `%s`", n)
				} else {
					messages.Infof("skipping file `%s`", n)
				}
				return yes
			},
		},
	})

	if err != nil {
		messages.Fatal(err)
		return
	}
}
