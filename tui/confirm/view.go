package confirm

import "github.com/razshare/frizzante/tui/config"

func (model *Model) View() string {
	if model.DefaultValue {
		return config.Styles.Menu.PaddingRight(1).Render("⎚") + config.Styles.Menu.Render(model.Prompt, "(Y/n)")
	}
	return config.Styles.Menu.PaddingRight(1).Render("⎚") + config.Styles.Menu.Render(model.Prompt, "(y/N)")
}
