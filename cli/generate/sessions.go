package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/singleselect"
)

func Sessions(options SessionsOptions) (err error) {
	if options.Type == "" {
		if options.Auto {
			options.Type = "memory"
		} else if options.Type, err = singleselect.Send([]search.Choice{{Id: "memory"}, {Id: "disk"}}, "type of sessions"); err != nil {
			return
		}
	}

	stype := strings.ToLower(options.Type)

	if !slices.Contains([]string{"memory", "disk"}, stype) {
		err = fmt.Errorf("sessions of type %s are not supported", stype)
		return
	}

	baseDirectory := filepath.Join("lib", "sessions", stype)

	if files.IsDirectory(baseDirectory) {
		if !options.Auto {
			var yes bool
			if yes, err = confirm.Sendf(true, "%s already exists. Overwrite?", baseDirectory); err != nil {
				return
			}

			if !yes {
				messages.Infof("skipping %s", baseDirectory)
				return
			}
		}

		if err = os.RemoveAll(baseDirectory); err != nil {
			return
		}
	}

	var from string
	if stype == "memory" {
		from = "internal/project/" + baseDirectory
	} else {
		from = "internal/additions/" + baseDirectory
	}

	if err = Copy(CopyOptions{
		From: from,
		To:   baseDirectory,
		Auto: options.Auto,
		Efs:  options.Efs,
	}); err != nil {
		return
	}

	if err = FixImports(FixImportsOptions{Directory: baseDirectory}); err != nil {
		return
	}

	switch stype {
	case "memory":
		messages.Success(
			"memory sessions generated.\n",
			baseDirectory+"/new.go\n",
			baseDirectory+"/start.go\n",
			baseDirectory+"/types.go\n",
		)
		messages.Tip(
			"## usage example\n",
			"func(client *clients.Client){\n",
			"    session := sessions.Start(receive.SessionId(client))\n",
			"}\n",
			"\n",
			"## state shape\n",
			"Your session state is defined by sessions.State,\n",
			"which is located in "+baseDirectory+"/types.go.\n",
			"\n",
			"## initial state\n",
			"Every new session is initialized with sessions.New(), \n",
			"which is located in "+baseDirectory+"/new.go.\n",
		)
	case "disk":
		messages.Success(
			"disk sessions generated at session.*\n",
			baseDirectory+"/new.go\n",
			baseDirectory+"/start.go\n",
			baseDirectory+"/types.go\n",
		)
		messages.Tip(
			"## usage example\n",
			"func(client *clients.Client){\n",
			"    session := sessions.Start(receive.SessionId(client))\n",
			"    defer sessions.Save(client, session)\n",
			"}\n",
			"\n",
			"## state shape\n",
			"session state is defined by sessions.State,\n",
			"which is located in "+baseDirectory+"/types.go.\n",
			"\n",
			"## initial state\n",
			"ever new session is initialized with sessions.New(), \n",
			"which is located in "+baseDirectory+"/new.go.\n",
		)
	}

	return
}
