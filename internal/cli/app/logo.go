package app

import (
	"github.com/razshare/frizzante/tui/config"
)

func Logo(options *App) (logo string, err error) {
	var data []byte
	if data, err = options.Efs.ReadFile("logo.txt"); err != nil {
		return
	}
	logo = config.Styles.BigText.Render(string(data))
	return
}
