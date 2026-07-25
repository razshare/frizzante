package menus

import (
	"github.com/razshare/frizzante/v2/cli/apps"
	"github.com/razshare/frizzante/v2/tui/search"
	"github.com/razshare/frizzante/v2/tui/select_one"
)

func Render(menu *Menu, app apps.App, value string, query []string, depth int) (id string, err error) {
	choices := make([]search.Choice, 0)
	for _, item := range menu.Items {
		if item.Hidden {
			continue
		}
		choices = append(choices, item.Choice)
	}
	if id, err = select_one.Send(choices, menu.Title); err != nil {
		return
	}
	for _, item := range menu.Items {
		if item.Choice.Id != id {
			continue
		}
		err = item.Handle(menu, app, value, query, depth)
		break
	}
	return
}
