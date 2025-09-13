package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/singleselect"
)

func Session(options SessionOptions) (err error) {
	if options.Type == "" {
		if options.Auto {
			options.Type = "memory"
		} else if options.Type, err = singleselect.Send([]search.Choice{{Id: "memory"}, {Id: "disk"}}, "session type"); err != nil {
			return
		}
	}

	if options.Type != "memory" && options.Type != "disk" {
		err = fmt.Errorf("session of type %s id not supported", options.Type)
		return
	}

	stype := strings.ToLower(options.Type)

	to := filepath.Join("lib", "session", stype)

	if files.IsDirectory(to) {
		if !options.Auto {
			var overwrite bool
			if overwrite, err = confirm.Sendf(true, "%s already exists. Overwrite?", to); err != nil {
				return
			}

			if !overwrite {
				messages.Infof("skipping %s", to)
				return
			}
		}

		if err = os.RemoveAll(to); err != nil {
			return
		}
	}

	var from string
	if stype == "memory" {
		from = "internal/project/" + to
	} else {
		from = "internal/additions/" + to
	}

	if err = Copy(CopyOptions{
		From: from,
		To:   to,
		Auto: options.Auto,
		Efs:  options.Efs,
	}); err != nil {
		return
	}

	if err = FixImports(FixImportsOptions{Directory: to}); err != nil {
		return
	}

	switch stype {
	case "memory":
		messages.Success(
			"memory session generated into session.*\n",
			to+"/new.go\n",
			to+"/start.go\n",
			to+"/types.go\n",
		)
		messages.Tip(
			"## usage example\n",
			"func(c *client.Client){\n",
			"    s := session.Start(receive.SessionId(c))\n",
			"}\n",
			"\n",
			"## state shape\n",
			"Your session state is defined by session.State,\n",
			"which is located in "+to+"/types.go.\n",
			"\n",
			"## initial state\n",
			"Every new session is initialized with session.New(), \n",
			"which is located in "+to+"/new.go.\n",
		)
	case "disk":
		messages.Success(
			"disk session generated at session.*\n",
			to+"/new.go\n",
			to+"/start.go\n",
			to+"/types.go\n",
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
			"which is located in "+to+"/types.go.\n",
			"\n",
			"## initial state\n",
			"ever new session is initialized with session.New(), \n",
			"which is located in "+to+"/new.go.\n",
		)
	}

	return
}
