package on

import (
	"embed"
	"fmt"
	"github.com/razshare/frizzante/cli/flags"
	"github.com/razshare/frizzante/tui/config"
	"github.com/razshare/frizzante/tui/input"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/singleselect"
)

func Menu(efs embed.FS) {
	options := []string{
		"Help",
		"Update",
		"Install",
		"Version",
		"Create Project",
		"Add",
		"Add?",
		"Test",
		"Package",
		"Package Watch",
		"Check",
		"Format",
		"Touch",
		"Clean",
		"Dev",
		"Build",
		"Configure",
		//"Create Sqlite Database",
		//"Sqlc Generate",
	}

	logo, readError := efs.ReadFile("clilogo.txt")
	if readError != nil {
		messages.Fatal(readError)
	}
	fmt.Println(config.Styles.BigText.Render(string(logo)))

	result := singleselect.Send(options, "Pick an option")

	if result == "Help" {
		*flags.Help = true
		Start(efs)
		return
	}

	if result == "Version" {
		*flags.Version = true
		Start(efs)
		return
	}

	if result == "Create Project" {
		*flags.CreateProject = input.Send("Give the project a name")
		Start(efs)
		return
	}

	if result == "Add" {
		*flags.Add = ":pick"
		Start(efs)
		return
	}

	if result == "Add?" {
		*flags.Add = "?"
		Start(efs)
		return
	}

	if result == "Test" {
		*flags.Test = true
		Start(efs)
		return
	}

	if result == "Package" {
		*flags.Package = true
		Start(efs)
		return
	}

	if result == "Package Watch" {
		*flags.PackageWatch = true
		Start(efs)
		return
	}

	if result == "Check" {
		*flags.Check = true
		Start(efs)
		return
	}

	if result == "Update" {
		*flags.Update = true
		Start(efs)
		return
	}

	if result == "Install" {
		*flags.Install = true
		Start(efs)
		return
	}

	if result == "Format" {
		*flags.Format = true
		Start(efs)
		return
	}

	if result == "Touch" {
		*flags.Touch = true
		Start(efs)
		return
	}

	if result == "Clean" {
		*flags.Clean = true
		Start(efs)
		return
	}

	if result == "Dev" {
		*flags.Dev = true
		Start(efs)
		return
	}

	if result == "Build" {
		*flags.Build = true
		Start(efs)
		return
	}

	if result == "Configure" {
		*flags.Configure = true
		Start(efs)
		return
	}

	if result == "Sqlc Generate" {
		*flags.SqlcGenerate = true
		Start(efs)
		return
	}

	if result == "Create Sqlite Database" {
		*flags.CreateSqliteDatabase = true
		Start(efs)
		return
	}
}
