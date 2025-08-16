package codegen

import (
	"github.com/razshare/frizzante/cli"
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/singleselect"
	"path/filepath"
	"strings"
)

func Session(c *cli.Cli, base string) error {
	tchoice, err := singleselect.Send(
		[]string{"memory", "disk"},
		"session type",
	)
	if err != nil {
		return err
	}

	t := strings.ToLower(tchoice)

	dst := filepath.Join(base, "lib", "session")

	err = Copy(c.Efs, []CopyInstruction{
		{
			From: "template/lib/session/" + t,
			To:   dst,
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
			dst+"/new.go\n",
			dst+"/start.go\n",
			dst+"/types.go\n",
		)
		messages.Tip(
			"## usage example\n",
			"func(c *client.Client){\n",
			"    s := session.Start(receive.SessionId(c))\n",
			"}\n",
			"\n",
			"## state shape\n",
			"Your session state is defined by session.State,\n",
			"which is located in "+dst+"/types.go.\n",
			"\n",
			"## initial state\n",
			"Every new session is initialized with session.New(), \n",
			"which is located in "+dst+"/new.go.\n",
		)
	case "disk":
		messages.Success(
			"disk session generated at session.*\n",
			dst+"/new.go\n",
			dst+"/start.go\n",
			dst+"/types.go\n",
		)
		messages.Tip(
			"## usage example\n",
			"func(c *client.Client){\n",
			"    s := session.Start(receive.SessionId(c))\n",
			"    defer session.Save(c, s)\n",
			"}\n",
			"\n",
			"## state shape\n",
			"session state is defined by session.State,\n",
			"which is located in "+dst+"/types.go.\n",
			"\n",
			"## initial state\n",
			"ever new session is initialized with session.New(), \n",
			"which is located in "+dst+"/new.go.\n",
		)
	}

	return nil
}
