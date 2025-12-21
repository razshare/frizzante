package cli

import (
	"github.com/razshare/frizzante/cli/menus"
)

func StartApp(options StartAppOptions) (err error) {
	app := options.App
	modifiers := app.Modifiers
	efs := app.Efs
	query := options.Query
	persistent := query == ""

	var menu *menus.Menu
	if menu, err = menus.NewMainMenu(menus.NewMainMenuOptions{
		Efs:        efs,
		Modifiers:  modifiers,
		Persistent: persistent,
	}); err != nil {
		return
	}

	err = menus.ParseQueryAndActivate(menu, query)

	return
}
