package generate

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/tui/confirm"
	messages2 "github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/singleselect"
)

func Session(options SessionOptions) (err error) {
	var choice string
	if choice, err = singleselect.Send([]search.Choice{{Id: "memory"}, {Id: "disk"}}, "session type"); err != nil {
		return
	}

	choice = strings.ToLower(choice)

	lib := filepath.Join("lib", "session", choice)

	if files.IsDirectory(lib) {
		if !options.Auto {
			var overwrite bool
			if overwrite, err = confirm.Sendf(true, "%s already exists. Overwrite?", lib); err != nil {
				return
			}

			if !overwrite {
				messages2.Infof("skipping %s", lib)
				return
			}
		}

		if err = os.RemoveAll(lib); err != nil {
			return
		}
	}

	if err = Copy(CopyOptions{
		From: "internal/project/lib/session/" + choice,
		To:   lib,
		Auto: options.Auto,
		Efs:  options.Efs,
	}); err != nil {
		return
	}

	switch choice {
	case "memory":
		messages2.Success(
			"memory session generated into session.*\n",
			lib+"/new.go\n",
			lib+"/start.go\n",
			lib+"/types.go\n",
		)
		messages2.Tip(
			"## usage example\n",
			"func(c *client.Client){\n",
			"    s := session.Start(receive.SessionId(c))\n",
			"}\n",
			"\n",
			"## state shape\n",
			"Your session state is defined by session.State,\n",
			"which is located in "+lib+"/types.go.\n",
			"\n",
			"## initial state\n",
			"Every new session is initialized with session.New(), \n",
			"which is located in "+lib+"/new.go.\n",
		)
	case "disk":
		messages2.Success(
			"disk session generated at session.*\n",
			lib+"/new.go\n",
			lib+"/start.go\n",
			lib+"/types.go\n",
		)
		messages2.Tip(
			"## usage example\n",
			"func(c *client.Client){\n",
			"    s := session.Start(receive.SessionId(c))\n",
			"    defer session.Save(c, s)\n",
			"}\n",
			"\n",
			"## state shape\n",
			"session state is defined by session.State,\n",
			"which is located in "+lib+"/types.go.\n",
			"\n",
			"## initial state\n",
			"ever new session is initialized with session.New(), \n",
			"which is located in "+lib+"/new.go.\n",
		)
	}

	return
}
