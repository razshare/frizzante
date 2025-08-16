package codegen

import (
	"embed"
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/singleselect"
	"path/filepath"
	"strings"
)

func Session(efs embed.FS, base string) error {
	tchoice, err := singleselect.Send(
		[]string{"Memory", "Disk"},
		"session type",
	)
	if err != nil {
		return err
	}

	t := strings.ToLower(tchoice)

	to := filepath.Join(base, "lib", "session")

	err = Copy(efs, []Generation{
		{
			From: "template/lib/session/" + t,
			To:   to,
			Overwrite: func(n string) (bool, error) {
				var overwrite bool

				overwrite, err = confirm.Sendf(true, "%s already exists. Overwrite?", n)
				if err != nil {
					return false, err
				}

				if overwrite {
					messages.Infof("overwriting %s", n)
				} else {
					messages.Infof("skipping %s", n)
				}

				return overwrite, nil
			},
		},
	})

	if err != nil {
		return err
	}

	switch t {
	case "memory":
		messages.Success(
			"memory session generated into session.*\n",
			to+"/new.go\n",
			to+"/start.go\n",
			to+"/types.go\n",
		)
		messages.Tip(
			"## Usage Example\n",
			"func(c *client.Client){\n",
			"    s := session.Start(receive.SessionId(c))\n",
			"}\n",
			"\n",
			"## State Shape\n",
			"Your session state is defined by session.State,\n",
			"which is located in "+to+"/types.go.\n",
			"\n",
			"## Initial State\n",
			"Every new session is initialized with session.New(), \n",
			"which is located in "+to+"/new.go.\n",
		)
	case "disk":
		messages.Success(
			"disk session generated into session.*\n",
			to+"/new.go\n",
			to+"/start.go\n",
			to+"/types.go\n",
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
			"which is located in "+to+"/types.go.\n",
			"\n",
			"## Initial State\n",
			"Every new session is initialized with session.New(), \n",
			"which is located in "+to+"/new.go.\n",
		)
	}

	return nil
}
