package menus

import (
	"errors"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/select_one"
)

func ParseQueryAndActivate(menu *Menu, query string) (err error) {
	if query != "" {
		for _, item := range menu.Items {
			for _, id := range item.Ids {
				if strings.HasPrefix(query, id) {
					value := strings.TrimSpace(strings.TrimPrefix(query, id))
					err = item.Handler(value)
					return
				}
			}
		}
	}

	if menu.Logo != "" {
		println(menu.Logo)
	}

	choices := make([]search.Choice, 0)

	for _, item := range menu.Items {
		if item.Hidden {
			continue
		}
		choices = append(choices, item.Choice)
	}

	for {
		var id string
		if id, err = select_one.Send(choices, menu.Title); err != nil {
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

			if err = item.Handler(""); err != nil {
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
