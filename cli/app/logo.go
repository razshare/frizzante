package app

import (
	"github.com/razshare/frizzante/tui/config"
)

func Logo(opts *App) (logo string, err error) {
	var data []byte
	if data, err = opts.Efs.ReadFile("clilogo.txt"); err != nil {
		return
	}
	logo = config.Styles.BigText.Render(string(data))
	return
}
