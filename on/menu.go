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
			Updates Go and JavaScript dependencies. This might bump version numbers.
		`,
		`Install
			Installs Go and JavaScript dependencies.
		`,
		`Version
			Shows the version number of this binary.
		`,
		`Create Project
			Creates a new Frizzante project. Once done, change directory into the project and run frizzante --configure or make configure.
		`,
		`Generate
			Generates code and resources. Select for more details.
		`,
		`Test
			Runs Go tests.
		`,
		`Package
			Packages the JavaScript/Svelte application into app/dist.
		`,
		`Package Watch
			Watches for changes in the JavaScript application and packages it automatically into app/dist.
		`,
		`Check
			Checks for JavaScript and Svelte and errors.
		`,
		`Format
			Formats Go, JavaScript and Svelte code.
		`,
		`Touch
			Touches the app/dist directory with placeholder files. This can be useful to silence //go:embed errors.
		`,
		`Clean
			Deletes Go temporary objects, app/dist, app/modules, .gen/tmp and .vite.
		`,
		`Dev
			Runs Air and Vite development server in parallel.
		`,
		`Build
			Builds the whole project into one single binary located at .gen/bin/app.
		`,
		`Configure
			Generates Air and Bun binaries then it installs Go and JavaScript dependencies.
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
