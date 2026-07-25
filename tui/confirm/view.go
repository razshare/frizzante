package confirm

import (
	"strings"

	"github.com/razshare/frizzante/v2/tui/configs"
)

func (model *Model) View() string {
	var builder strings.Builder
	builder.WriteString(configs.Styles.Menu.Render("⎚"))
	builder.WriteString(configs.Styles.Menu.PaddingLeft(1).PaddingRight(1).Render(model.Prompt))
	if model.Confirmed {
		builder.WriteString(configs.Styles.Menu.Render("● Yes"))
		builder.WriteString(configs.Styles.UserGuide.PaddingLeft(1).PaddingRight(1).Render("/"))
		builder.WriteString(configs.Styles.UserGuide.Render("○ No"))
	} else {
		builder.WriteString(configs.Styles.UserGuide.Render("○ Yes"))
		builder.WriteString(configs.Styles.UserGuide.PaddingLeft(1).PaddingRight(1).Render("/"))
		builder.WriteString(configs.Styles.Menu.Render("● No"))
	}
	return builder.String()
}
