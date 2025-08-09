package cli

import (
	"embed"
	"github.com/razshare/frizzante/tui/input"
	"github.com/razshare/frizzante/tui/singleselect"
)

func OnMenu(efs embed.FS) {
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
		"Sqlc Generate",
	}

	result, showError := singleselect.Send("Pick an option", options)
	if showError != nil {
		Fatal(showError)
	}

	if result == "Help" {
		*FlagHelp = true
		OnStart(efs)
		return
	}

	if result == "Version" {
		*FlagVersion = true
		OnStart(efs)
		return
	}

	if result == "Create Project" {
		projectName, projectNameError := input.Send("Give the project a name")
		if projectNameError != nil {
			Fatal(projectNameError)
		}
		*FlagCreateProject = projectName
		OnStart(efs)
		return
	}

	if result == "Add" {
		*FlagAdd = ":pick"
		OnStart(efs)
		return
	}

	if result == "Add?" {
		*FlagAdd = "?"
		OnStart(efs)
		return
	}

	if result == "Test" {
		*FlagTest = true
		OnStart(efs)
		return
	}

	if result == "Package" {
		*FlagPackage = true
		OnStart(efs)
		return
	}

	if result == "Package Watch" {
		*FlagPackageWatch = true
		OnStart(efs)
		return
	}

	if result == "Check" {
		*FlagCheck = true
		OnStart(efs)
		return
	}

	if result == "Update" {
		*FlagUpdate = true
		OnStart(efs)
		return
	}

	if result == "Install" {
		*FlagInstall = true
		OnStart(efs)
		return
	}

	if result == "Format" {
		*FlagFormat = true
		OnStart(efs)
		return
	}

	if result == "Touch" {
		*FlagTouch = true
		OnStart(efs)
		return
	}

	if result == "Clean" {
		*FlagClean = true
		OnStart(efs)
		return
	}

	if result == "Dev" {
		*FlagDev = true
		OnStart(efs)
		return
	}

	if result == "Build" {
		*FlagBuild = true
		OnStart(efs)
		return
	}

	if result == "Configure" {
		*FlagConfigure = true
		OnStart(efs)
		return
	}

	if result == "Sqlc Generate" {
		*FlagSqlcGenerate = true
		OnStart(efs)
		return
	}
}
