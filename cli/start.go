package cli

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/cli/apps"
	"github.com/razshare/frizzante/cli/menus"
	"github.com/razshare/frizzante/tui/config"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/select_one"
)

func Start(app *apps.App) (err error) {
	messages.Prefix = config.Styles.Menu.PaddingRight(1).Render("│")
	var menu *menus.Menu
	menu, err = menus.New(app)
	if err != nil {
		return
	}

	choices := make([]search.Choice, 0)

	for _, item := range menu.Items {
		if item.Hidden {
			continue
		}
		choices = append(choices, item.Choice)
	}

	// If this for loop returns,
	// it means the choice has been inlined.
	for _, item := range menu.Items {
		if !item.Active() {
			continue
		}

		return item.Handler()
	}

	var logo string
	if logo, err = apps.Logo(app); err != nil {
		messages.Error(err)
		err = nil
	} else {
		println(logo)
	}

	// If we reach this point,
	// it means we need to show the TUI menu.
	for {
		var id string
		if id, err = select_one.Send(choices, "menu"); err != nil {
			if errors.Is(err, tea.ErrInterrupted) {
				return
			}
			messages.Error(err)
			err = nil
			continue
		}

		for _, item := range menu.Items {
			if item.Choice.Id != id {
				continue
			}

			if err = item.Handler(); err != nil {
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
