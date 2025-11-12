package apps

import (
	"github.com/razshare/frizzante/tui/configs"
)

func Logo(options *App) (logo string, err error) {
	var data []byte
	if data, err = options.Efs.ReadFile("logo.txt"); err != nil {
		return
	}
	logo = configs.Styles.BigText.PaddingLeft(1).PaddingRight(1).Render(string(data))
	return
}
