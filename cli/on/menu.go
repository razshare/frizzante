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

var Functions = map[string]func(c *cli.Cli) (func(), error){
	"configure": func(c *cli.Cli) (func(), error) {
		*c.Flags.Configure = true
		return func() {
			*c.Flags.Configure = false
		}, nil
	},
	"create project": func(c *cli.Cli) (func(), error) {
		var err error
		*c.Flags.CreateProject, err = input.Send("give the project a name")
		if err != nil {
			return nil, err
		}
		return func() {
			*c.Flags.CreateProject = ""
		}, nil
	},
	"dev": func(c *cli.Cli) (func(), error) {
		*c.Flags.Dev = true
		return func() {
			*c.Flags.Dev = false
		}, nil
	},
	"build": func(c *cli.Cli) (func(), error) {
		*c.Flags.Build = true
		return func() {
			*c.Flags.Build = false
		}, nil
	},
	"install": func(c *cli.Cli) (func(), error) {
		*c.Flags.Install = true
		return func() {
			*c.Flags.Install = false
		}, nil
	},
	"update": func(c *cli.Cli) (func(), error) {
		*c.Flags.Update = true
		return func() {
			*c.Flags.Update = false
		}, nil
	},
	"generate": func(c *cli.Cli) (func(), error) {
		*c.Flags.Generate = ":pick"
		return func() {
			*c.Flags.Generate = ""
		}, nil
	},
	"package": func(c *cli.Cli) (func(), error) {
		*c.Flags.Package = true
		return func() {
			*c.Flags.Package = false
		}, nil
	},
	"package (watch)": func(c *cli.Cli) (func(), error) {
		*c.Flags.PackageWatch = true
		return func() {
			*c.Flags.PackageWatch = false
		}, nil
	},
	"test": func(c *cli.Cli) (func(), error) {
		*c.Flags.Test = true
		return func() {
			*c.Flags.Test = false
		}, nil
	},
	"check": func(c *cli.Cli) (func(), error) {
		*c.Flags.Check = true
		return func() {
			*c.Flags.Check = false
		}, nil
	},
	"format": func(c *cli.Cli) (func(), error) {
		*c.Flags.Format = true
		return func() {
			*c.Flags.Format = false
		}, nil
	},
	"touch": func(c *cli.Cli) (func(), error) {
		*c.Flags.Touch = true
		return func() {
			*c.Flags.Touch = false
		}, nil
	},
	"clean": func(c *cli.Cli) (func(), error) {
		*c.Flags.Clean = true
		return func() {
			*c.Flags.Clean = false
		}, nil
	},
	"help": func(c *cli.Cli) (func(), error) {
		*c.Flags.Help = true
		return func() {
			*c.Flags.Help = false
		}, nil
	},
	"version": func(c *cli.Cli) (func(), error) {
		*c.Flags.Version = true
		return func() {
			*c.Flags.Version = false
		}, nil
	},
}

func Menu(c *cli.Cli, _ bool, base string) error {
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
			builds ` + *c.Flags.App + `
		`,
		`package (watch)
			builds ` + *c.Flags.App + ` on change
		`,
		`check
			checks for code errors
		`,
		`format
			formats code
		`,
		`touch
			adds placeholders in ` + *c.Flags.App + `/dist
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

		revert, err = choose(c)
		if err != nil {
			messages.Error(err)
		}

		err = Start(c, true, base)
		if err != nil {
			if errors.Is(err, tea.ErrInterrupted) {
				os.Exit(0)
			}

			messages.Error(err)
		}

		revert()

		err = Menu(c, true, base)
		if err != nil {
			return err
		}
	}

	return nil
}
