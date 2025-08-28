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

func Session(options SessionOptions) (err error) {
	if files.IsDirectory(options.Lib) {
		if !options.Auto {
			var overwrite bool
			if overwrite, err = confirm.Sendf(true, "%s already exists. Overwrite?", options.Lib); err != nil {
				return
			}

			if !overwrite {
				messages.Infof("skipping %s", options.Lib)
				return
			}
		}

		if err = os.RemoveAll(options.Lib); err != nil {
			return
		}
	}

	var choice string
	if choice, err = singleselect.Send([]search.Choice{{Id: "memory"}, {Id: "disk"}}, "session type"); err != nil {
		return
	}

	choice = strings.ToLower(choice)

	if err = Copy(CopyOptions{From: "internal/template/lib/session/" + choice, To: options.Lib, Auto: options.Auto}); err != nil {
		return
	}

	switch choice {
	case "memory":
		messages.Success(
			"memory session generated into session.*\n",
			options.Lib+"/new.go\n",
			options.Lib+"/start.go\n",
			options.Lib+"/types.go\n",
		)
		messages.Tip(
			"## usage example\n",
			"func(c *client.Client){\n",
			"    s := session.Start(receive.SessionId(c))\n",
			"}\n",
			"\n",
			"## state shape\n",
			"Your session state is defined by session.State,\n",
			"which is located in "+options.Lib+"/types.go.\n",
			"\n",
			"## initial state\n",
			"Every new session is initialized with session.New(), \n",
			"which is located in "+options.Lib+"/new.go.\n",
		)
	case "disk":
		messages.Success(
			"disk session generated at session.*\n",
			options.Lib+"/new.go\n",
			options.Lib+"/start.go\n",
			options.Lib+"/types.go\n",
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
			"which is located in "+options.Lib+"/types.go.\n",
			"\n",
			"## initial state\n",
			"ever new session is initialized with session.New(), \n",
			"which is located in "+options.Lib+"/new.go.\n",
		)
	}

	return
}
