package menus

import "github.com/razshare/frizzante/cli/apps"

func Activate(menu *Menu, app apps.App, args []string) (id string, err error) {
	argsLength := len(args)
	var value string
	var query []string
	if argsLength == 0 {
		query = make([]string, 0)
	} else if argsLength == 1 {
		value = args[0]
		query = make([]string, 0)
	} else if argsLength >= 2 {
		value = args[0]
		query = args[1:]
	}
	for _, item := range menu.Items {
		if !item.Active(menu, app, value, query) {
			continue
		}
		queryLength := len(query)
		id = item.Choice.Id
		var valueLocal string
		var queryLocal []string
		if queryLength == 0 {
			queryLocal = make([]string, 0)
		} else if queryLength == 1 {
			valueLocal = query[0]
			queryLocal = make([]string, 0)
		} else {
			valueLocal = query[0]
			queryLocal = query[1:]
		}
		err = item.Handle(menu, app, valueLocal, queryLocal)
		break
	}
	return
}
