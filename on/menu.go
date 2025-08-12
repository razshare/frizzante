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
			Deletes removes Go temporary objects, deletes app/dist, app/modules, .gen/tmp and .vite.
		`,
		`Dev
			Runs Air and Vite development server in parallel.
		`,
		`Build
			Build the whole project into one single binary located at .gen/bin/app.
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

	if result == "Generate" {
		*flags.Generate = ":pick"
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
}
