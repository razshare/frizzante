package cli

import (
	"errors"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/cli/app"
	_menu "github.com/razshare/frizzante/cli/menu"
	"github.com/razshare/frizzante/tui/config"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/singleselect"
)

func Start(a *app.App) (err error) {
	var menu *_menu.Menu
	menu, err = _menu.New(a)
	if err != nil {
		return
	}

	chs := make([]search.Choice, 0)

	for _, it := range menu.Items {
		if it.Hidden {
			continue
		}
		chs = append(chs, it.Choice)
	}

	// If this for loop returns,
	// it means the choice has been inlined.
	for _, it := range menu.Items {
		if !it.Active() {
			continue
		}

		print(config.Styles.Menu.PaddingRight(1).Render("⎚"))
		println(config.Styles.Menu.Render(fmt.Sprintf("running ▷ %s (%s)", it.Choice.Id, it.Choice.Description)))
		messages.Prefix = config.Styles.Menu.PaddingRight(1).Render("│")

		return it.Handler()
	}

	var logo string
	if logo, err = app.Logo(a); err != nil {
		messages.Error(err)
		err = nil
	} else {
		println(logo)
	}

	print(config.Styles.Menu.PaddingRight(1).Render("⎚"))
	println(config.Styles.Menu.Render("menu"))
	messages.Prefix = config.Styles.Menu.PaddingRight(1).Render("│")

	// If we reach this point,
	// it means we need to show the TUI menu.
	for {
		var id string
		id, err = singleselect.Send(chs, "menu")
		if err != nil {
			if errors.Is(err, tea.ErrInterrupted) {
				return
			}
			messages.Error(err)
			err = nil
			continue
		}

		for _, it := range menu.Items {
			if it.Choice.Id != id {
				continue
			}

			if err = it.Handler(); err != nil {
				if errors.Is(err, tea.ErrInterrupted) {
					return
				}
				messages.Error(err)
				err = nil
			}
			break
		}
	}
}
