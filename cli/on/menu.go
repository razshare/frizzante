package on

import (
	"errors"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/cli"
	"github.com/razshare/frizzante/tui/input"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/singleselect"
	"os"
)

var Functions = map[string]func() (func(), error){
	"configure": func() (func(), error) {
		*cli.Configure = true
		return func() {
			*cli.Configure = false
		}, nil
	},
	"create project": func() (func(), error) {
		var err error
		*cli.CreateProject, err = input.Send("give the project a name")
		if err != nil {
			return nil, err
		}
		return func() {
			*cli.CreateProject = ""
		}, nil
	},
	"dev": func() (func(), error) {
		*cli.Dev = true
		return func() {
			*cli.Dev = false
		}, nil
	},
	"build": func() (func(), error) {
		*cli.Build = true
		return func() {
			*cli.Build = false
		}, nil
	},
	"install": func() (func(), error) {
		*cli.Install = true
		return func() {
			*cli.Install = false
		}, nil
	},
	"update": func() (func(), error) {
		*cli.Update = true
		return func() {
			*cli.Update = false
		}, nil
	},
	"generate": func() (func(), error) {
		*cli.Generate = ":pick"
		return func() {
			*cli.Generate = ""
		}, nil
	},
	"package": func() (func(), error) {
		*cli.Package = true
		return func() {
			*cli.Package = false
		}, nil
	},
	"package (watch)": func() (func(), error) {
		*cli.PackageWatch = true
		return func() {
			*cli.PackageWatch = false
		}, nil
	},
	"test": func() (func(), error) {
		*cli.Test = true
		return func() {
			*cli.Test = false
		}, nil
	},
	"check": func() (func(), error) {
		*cli.Check = true
		return func() {
			*cli.Check = false
		}, nil
	},
	"format": func() (func(), error) {
		*cli.Format = true
		return func() {
			*cli.Format = false
		}, nil
	},
	"touch": func() (func(), error) {
		*cli.Touch = true
		return func() {
			*cli.Touch = false
		}, nil
	},
	"clean": func() (func(), error) {
		*cli.Clean = true
		return func() {
			*cli.Clean = false
		}, nil
	},
	"help": func() (func(), error) {
		*cli.Help = true
		return func() {
			*cli.Help = false
		}, nil
	},
	"version": func() (func(), error) {
		*cli.Version = true
		return func() {
			*cli.Version = false
		}, nil
	},
}

func Menu() error {
	options := []string{
		`configure
			installs required binaries and dependencies
		`,
		`create project
			creates a new project and configures it
		`,
		`update
			updates dependencies
		`,
		`install
			installs dependencies
		`,
		`version
			shows binary version
		`,
		`generate
			generates code and resources
		`,
		`test
			runs tests
		`,
		`package
			builds ` + "app" + `
		`,
		`package (watch)
			builds ` + "app" + ` on change
		`,
		`check
			checks for code errors
		`,
		`format
			formats code
		`,
		`touch
			adds placeholders in ` + "app" + `/dist
		`,
		`clean
			deletes unnecessary files
		`,
		`dev
			runs Air and Vite in parallel
		`,
		`build
			builds project at .gen/bin/app.
		`,
	}

	choice, err := singleselect.Send(options, "menu")

	if err != nil {
		return err
	}

	if choose := Functions[choice]; choose != nil {
		var revert func()

		revert, err = choose()
		if err != nil {
			messages.Error(err)
		}

		err = Start()
		if err != nil {
			if errors.Is(err, tea.ErrInterrupted) {
				os.Exit(0)
			}

			messages.Error(err)
		}

		revert()

		err = Menu()
		if err != nil {
			return err
		}
	}

	return nil
}
