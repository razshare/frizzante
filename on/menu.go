package on

import (
	"embed"
	"fmt"
	"github.com/razshare/frizzante/cli/state"
	"github.com/razshare/frizzante/tui/config"
	"github.com/razshare/frizzante/tui/input"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/singleselect"
)

func Menu(efs embed.FS) {
	options := []string{
		`Help`,
		`Update
			Updates dependencies.
		`,
		`Install
			Installs dependencies.
		`,
		`Version
			Shows the version number of this binary.
		`,
		`Create Project
			Creates a new Frizzante project.
		`,
		`Generate
			Generates code and resources.
		`,
		`Test
			Runs Go tests.
		`,
		`Package
			Packages the app directory.
		`,
		`Package Watch
			Watches for changes and packages the app directory.
		`,
		`Check
			Checks for JavaScript and Svelte code errors.
		`,
		`Format
			Formats Go, JavaScript and Svelte code.
		`,
		`Touch
			Adds placeholder files in app/dist. This can be useful to silence //go:embed errors.
		`,
		`Clean
			Deletes Go temporary objects, app/dist, app/modules, .gen/tmp and .vite.
		`,
		`Dev
			Runs Air and Vite in parallel.
		`,
		`Build
			Builds project into a binary located at .gen/bin/app.
		`,
		`Configure
			Installs required binaries in .gen/bin and installs code dependencies.
		`,
	}

	logo, err := efs.ReadFile("clilogo.txt")
	if err != nil {
		messages.Fatal(err)
	}
	fmt.Println(config.Styles.BigText.Render(string(logo)))

	result := singleselect.Send(options, "Pick an option")

	if result == "Help" {
		*state.Help = true
		Start(efs)
		return
	}

	if result == "Version" {
		*state.Version = true
		Start(efs)
		return
	}

	if result == "Create Project" {
		*state.CreateProject = input.Send("Give the project a name")
		Start(efs)
		return
	}

	if result == "Generate" {
		*state.Generate = ":pick"
		Start(efs)
		return
	}

	if result == "Test" {
		*state.Test = true
		Start(efs)
		return
	}

	if result == "Package" {
		*state.Package = true
		Start(efs)
		return
	}

	if result == "Package Watch" {
		*state.PackageWatch = true
		Start(efs)
		return
	}

	if result == "Check" {
		*state.Check = true
		Start(efs)
		return
	}

	if result == "Update" {
		*state.Update = true
		Start(efs)
		return
	}

	if result == "Install" {
		*state.Install = true
		Start(efs)
		return
	}

	if result == "Format" {
		*state.Format = true
		Start(efs)
		return
	}

	if result == "Touch" {
		*state.Touch = true
		Start(efs)
		return
	}

	if result == "Clean" {
		*state.Clean = true
		Start(efs)
		return
	}

	if result == "Dev" {
		*state.Dev = true
		Start(efs)
		return
	}

	if result == "Build" {
		*state.Build = true
		Start(efs)
		return
	}

	if result == "Configure" {
		*state.Configure = true
		Start(efs)
		return
	}
}
