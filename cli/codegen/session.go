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

func Session(opts SessionOptions) (err error) {
	if files.IsDirectory(opts.Lib) {
		if !opts.Auto {
			var overwrite bool
			if overwrite, err = confirm.Sendf(true, "%s already exists. Overwrite?", opts.Lib); err != nil {
				return
			}

			if !overwrite {
				messages.Infof("skipping %s", opts.Lib)
				return
			}
		}

		if err = os.RemoveAll(opts.Lib); err != nil {
			return
		}
	}

	var choice string
	if choice, err = singleselect.Send([]search.Choice{{Id: "memory"}, {Id: "disk"}}, "session type"); err != nil {
		return
	}

	choice = strings.ToLower(choice)

	if err = Copy(CopyOptions{From: "internal/template/lib/session/" + choice, To: opts.Lib, Auto: opts.Auto}); err != nil {
		return
	}

	switch choice {
	case "memory":
		messages.Success(
			"memory session generated into session.*\n",
			opts.Lib+"/new.go\n",
			opts.Lib+"/start.go\n",
			opts.Lib+"/types.go\n",
		)
		messages.Tip(
			"## usage example\n",
			"func(c *client.Client){\n",
			"    s := session.Start(receive.SessionId(c))\n",
			"}\n",
			"\n",
			"## state shape\n",
			"Your session state is defined by session.State,\n",
			"which is located in "+opts.Lib+"/types.go.\n",
			"\n",
			"## initial state\n",
			"Every new session is initialized with session.New(), \n",
			"which is located in "+opts.Lib+"/new.go.\n",
		)
	case "disk":
		messages.Success(
			"disk session generated at session.*\n",
			opts.Lib+"/new.go\n",
			opts.Lib+"/start.go\n",
			opts.Lib+"/types.go\n",
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
			"which is located in "+opts.Lib+"/types.go.\n",
			"\n",
			"## initial state\n",
			"ever new session is initialized with session.New(), \n",
			"which is located in "+opts.Lib+"/new.go.\n",
		)
	}

	return
}
