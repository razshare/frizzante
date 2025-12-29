package cli

import (
	"github.com/razshare/frizzante/cli/menus"
	"github.com/razshare/frizzante/tui/configs"
	"github.com/razshare/frizzante/tui/messages"
)

func Start(options StartOptions) (err error) {
	app := options.App
	messages.Prefix = configs.Styles.Menu.PaddingRight(1).Render("│")
	_, err = menus.Activate(&menus.Main, app)
	return
}
