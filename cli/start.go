package cli

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	app2 "github.com/razshare/frizzante/cli/app"
	"github.com/razshare/frizzante/cli/menu"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/singleselect"
)

func Start(a *app2.App) error {
	m, err := menu.New(a)
	if err != nil {
		return err
	}

	chs := make([]search.Choice, 0)

	for _, it := range m.Items {
		if it.Hidden {
			continue
		}
		chs = append(chs, it.Choice)
	}

	// If this for loop returns,
	// it means the choice has been inlined.
	for _, it := range m.Items {
		if !it.Active() {
			continue
		}

		return it.Handler()
	}

	var logo string
	if logo, err = app2.Logo(a); err == nil {
		println(logo)
	}

	// If we reach this point,
	// it means we need to show the TUI menu.
	for {
		var id string
		id, err = singleselect.Send(chs, "menu")
		if err != nil {
			if errors.Is(err, tea.ErrInterrupted) {
				return err
			}
			messages.Error(err)
			continue
		}

		for _, it := range m.Items {
			if it.Choice.Id != id {
				continue
			}

			if err = it.Handler(); err != nil {
				if errors.Is(err, tea.ErrInterrupted) {
					return err
				}
				messages.Error(err)
			}
			break
		}
	}
}
