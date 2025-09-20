package cli

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/cli/app"
	_menu "github.com/razshare/frizzante/cli/menu"
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

		return it.Handler()
	}

	var logo string
	if logo, err = app.Logo(a); err != nil {
		messages.Error(err)
		err = nil
	} else {
		println(logo)
	}

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
