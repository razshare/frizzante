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
	"github.com/razshare/frizzante/tui/select_one"
)

func Sessions(options SessionsOptions) (err error) {
	if options.Type == "" {
		if options.Auto {
			options.Type = "memory"
		} else if options.Type, err = select_one.Send([]search.Choice{{Id: "memory"}, {Id: "disk"}}, "type of sessions"); err != nil {
			return
		}
	}

	sessionType := strings.ToLower(options.Type)

	if !slices.Contains([]string{"memory", "disk"}, sessionType) {
		err = fmt.Errorf("sessions of type %s are not supported", sessionType)
		return
	}

	baseDirectory := filepath.Join("lib", sessionType, "sessions")

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
	if sessionType == "memory" {
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

	switch sessionType {
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
			"    session := sessions.Start(client)\n",
			"}\n",
			"\n",
			"## session shape\n",
			"session shape is defined in "+baseDirectory+"/types.go.\n",
			"\n",
			"## initial state\n",
			"every new session is initialized with sessions.New(), \n",
			"which is located in "+baseDirectory+"/new.go.\n",
		)
	case "disk":
		messages.Success(
			"disk sessions generated.\n",
			baseDirectory+"/new.go\n",
			baseDirectory+"/start.go\n",
			baseDirectory+"/types.go\n",
		)
		messages.Tip(
			"## usage example\n",
			"func(client *clients.Client){\n",
			"    session := sessions.Start(client)\n",
			"    defer sessions.Start(session, client)\n",
			"}\n",
			"\n",
			"## session shape\n",
			"session shape is defined in "+baseDirectory+"/types.go.\n",
			"\n",
			"## initial state\n",
			"every new session is initialized with sessions.New(), \n",
			"which is located in "+baseDirectory+"/new.go.\n",
		)
	}

	return
}
