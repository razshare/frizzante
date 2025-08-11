package codegen

import (
	"embed"
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/singleselect"
	"path/filepath"
	"strings"
)

func Session(efs embed.FS) {
	t := strings.ToLower(singleselect.Send(
		[]string{"Memory", "Disk"},
		"How should the session be managed?",
	))

	err := Generate(efs, []Generation{
		{
			From: "template/lib/session/" + t,
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

	switch t {
	case "memory":
		messages.Success(
			"memory session generated into session.*\n",
			"./lib/session/new.go\n",
			"./lib/session/start.go\n",
			"./lib/session/types.go\n",
		)
		messages.Tip(
			"## Usage Example\n",
			"func(c *client.Client){\n",
			"    s := session.Start(receive.SessionId(c))\n",
			"}\n",
			"\n",
			"## State Shape\n",
			"Your session state is defined by session.State,\n",
			"which is located in ./lib/session/types.go.\n",
			"\n",
			"## Initial State\n",
			"Every new session is initialized with session.New(), \n",
			"which is located in ./lib/session/new.go.\n",
		)
	case "disk":
		messages.Success(
			"disk session generated into session.*\n",
			"./lib/session/new.go\n",
			"./lib/session/start.go\n",
			"./lib/session/types.go\n",
		)
		messages.Tip(
			"## Usage Example\n",
			"func(c *client.Client){\n",
			"    s := session.Start(receive.SessionId(c))\n",
			"    defer session.Save(c, s)\n",
			"}\n",
			"\n",
			"## State Shape\n",
			"Your session state is defined by session.State,\n",
			"which is located in ./lib/session/types.go.\n",
			"\n",
			"## Initial State\n",
			"Every new session is initialized with session.New(), \n",
			"which is located in ./lib/session/new.go.\n",
		)
	}
}
