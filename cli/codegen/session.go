package codegen

import (
	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/singleselect"
	"os"
	"strings"
)

func Session(o SessionOptions) error {
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

	tchoice, err := singleselect.Send(
		[]search.Choice{{Id: "memory"}, {Id: "disk"}},
		"session type",
	)

	if err != nil {
		return err
	}

	t := strings.ToLower(tchoice)

	err = Copy([]CopyInstruction{
		{
			Efs:  o.Efs,
			From: "template/lib/session/" + t,
			To:   o.Lib,
			Overwrite: func(n string) (bool, error) {
				if o.Auto {
					return true, nil
				}

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
			o.Lib+"/new.go\n",
			o.Lib+"/start.go\n",
			o.Lib+"/types.go\n",
		)
		messages.Tip(
			"## usage example\n",
			"func(c *client.Client){\n",
			"    s := session.Start(receive.SessionId(c))\n",
			"}\n",
			"\n",
			"## state shape\n",
			"Your session state is defined by session.State,\n",
			"which is located in "+o.Lib+"/types.go.\n",
			"\n",
			"## initial state\n",
			"Every new session is initialized with session.New(), \n",
			"which is located in "+o.Lib+"/new.go.\n",
		)
	case "disk":
		messages.Success(
			"disk session generated at session.*\n",
			o.Lib+"/new.go\n",
			o.Lib+"/start.go\n",
			o.Lib+"/types.go\n",
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
			"which is located in "+o.Lib+"/types.go.\n",
			"\n",
			"## initial state\n",
			"ever new session is initialized with session.New(), \n",
			"which is located in "+o.Lib+"/new.go.\n",
		)
	}

	return nil
}
