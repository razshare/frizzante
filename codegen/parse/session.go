package parse

import (
	"embed"
	"github.com/razshare/frizzante/codegen"
	"github.com/razshare/frizzante/embeds"
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/singleselect"
	"path/filepath"
	"strings"
)

func Session(efs embed.FS) {
	t := singleselect.Send(
		[]string{
			"Memory",
			"Disk",
		},
		"How should the session be managed?",
	)

	err := embeds.Generate(efs, []codegen.Generation{
		{
			From: "template/lib/session/" + strings.ReplaceAll(strings.ToLower(t), " ", ""),
			To:   filepath.Join("lib", "session"),
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
