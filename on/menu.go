package on

import (
	"embed"
	"fmt"
	"github.com/razshare/frizzante/cli/state"
	"github.com/razshare/frizzante/tui/config"
	"github.com/razshare/frizzante/tui/input"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/singleselect"
	"github.com/razshare/frizzante/tui/text"
)

var Functions = map[string]func() (func(), error){
	"configure": func() (func(), error) {
		*state.Configure = true
		return func() {
			*state.Configure = false
		}, nil
	},
	"create project": func() (func(), error) {
		var err error
		*state.CreateProject, err = input.Send("give the project a name")
		if err != nil {
			return nil, err
		}
		return func() {
			*state.CreateProject = ""
		}, nil
	},
	"dev": func() (func(), error) {
		*state.Dev = true
		return func() {
			*state.Dev = false
		}, nil
	},
	"build": func() (func(), error) {
		*state.Build = true
		return func() {
			*state.Build = false
		}, nil
	},
	"install": func() (func(), error) {
		*state.Install = true
		return func() {
			*state.Install = false
		}, nil
	},
	"update": func() (func(), error) {
		*state.Update = true
		return func() {
			*state.Update = false
		}, nil
	},
	"generate": func() (func(), error) {
		*state.Generate = ":pick"
		return func() {
			*state.Generate = ""
		}, nil
	},
	"package": func() (func(), error) {
		*state.Package = true
		return func() {
			*state.Package = false
		}, nil
	},
	"package (watch)": func() (func(), error) {
		*state.PackageWatch = true
		return func() {
			*state.PackageWatch = false
		}, nil
	},
	"test": func() (func(), error) {
		*state.Test = true
		return func() {
			*state.Test = false
		}, nil
	},
	"check": func() (func(), error) {
		*state.Check = true
		return func() {
			*state.Check = false
		}, nil
	},
	"format": func() (func(), error) {
		*state.Format = true
		return func() {
			*state.Format = false
		}, nil
	},
	"touch": func() (func(), error) {
		*state.Touch = true
		return func() {
			*state.Touch = false
		}, nil
	},
	"clean": func() (func(), error) {
		*state.Clean = true
		return func() {
			*state.Clean = false
		}, nil
	},
	"help": func() (func(), error) {
		*state.Help = true
		return func() {
			*state.Help = false
		}, nil
	},
	"version": func() (func(), error) {
		*state.Version = true
		return func() {
			*state.Version = false
		}, nil
	},
}

func Menu(efs embed.FS, clear bool) error {
	if !clear {
		logo, err := efs.ReadFile("clilogo.txt")
		if err != nil {
			return err
		}
		fmt.Println(config.Styles.BigText.Render(string(logo)))
	}

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
			builds ` + *state.App + `
		`,
		`package (watch)
			builds ` + *state.App + ` on change
		`,
		`check
			checks for code errors
		`,
		`format
			formats code
		`,
		`touch
			adds placeholders in ` + *state.App + `/dist
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

	if enactChoice := Functions[choice]; enactChoice != nil {
		var revert func()

		revert, err = enactChoice()
		if err != nil {
			messages.Error(err)
		}

		if clear {
			text.Clrscr()
		}

		err = Start(efs)
		if err != nil {
			messages.Error(err)
		}

		revert()

		err = Menu(efs, true)
		if err != nil {
			return err
		}
	}

	return nil
}
