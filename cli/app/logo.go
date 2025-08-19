package app

import (
	"github.com/razshare/frizzante/tui/config"
)

func Logo(o *App) (string, error) {
	var d []byte
	d, err := o.Efs.ReadFile("clilogo.txt")
	if err != nil {
		return "", err
	}

	return config.Styles.BigText.Render(string(d)), nil
}
