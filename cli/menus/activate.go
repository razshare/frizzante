package menus

import "github.com/razshare/frizzante/cli/apps"

func Activate(menu *Menu, app apps.App) (id string, err error) {
	for _, item := range menu.Items {
		if !item.Active(menu, app) {
			continue
		}
		id = item.Choice.Id
		err = item.Handle(menu, app)
		break
	}
	return
}
