package menus

import (
	"github.com/razshare/frizzante/v2/cli/apps"
	"github.com/razshare/frizzante/v2/tui/search"
)

type Item struct {
	Hidden bool
	Choice search.Choice
	Active func(menu *Menu, app apps.App, value string, query []string) (active bool)
	Handle func(menu *Menu, app apps.App, value string, query []string, depth int) (err error)
}
