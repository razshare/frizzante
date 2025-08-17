package cli

import (
	"errors"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/singleselect"
	"os"
)

func Select(c *Cli, o SelectOptions) error {
	options := []string{
		`configure
			installs required binaries and dependencies
		`,
		`create project
			creates a new project
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
		`reset
			deletes ~/.frizzante
		`,
	}

	choice, err := singleselect.Send(options, "menu")

	if err != nil {
		return err
	}

	if choose := c.Menu[choice]; choose != nil {
		var revert func()

		revert, err = choose()
		if err != nil {
			messages.Error(err)
		}

		err = Next(c, NextOptions{
			Go:       o.Go,
			Air:      o.Air,
			Bun:      o.Bun,
			Sqlc:     o.Sqlc,
			Platform: o.Platform,
			Action:   o.Action,
			Counter:  o.NextCounter,
		})

		if err != nil {
			if !errors.Is(err, tea.ErrInterrupted) {
				messages.Error(err)
			}
			os.Exit(0)
		}

		revert()

		err = Next(c, NextOptions{
			Go:       o.Go,
			Air:      o.Air,
			Bun:      o.Bun,
			Sqlc:     o.Sqlc,
			Platform: o.Platform,
			Action:   o.Action,
			Counter:  o.NextCounter,
		})

		if err != nil {
			return err
		}
	}

	return nil
}
