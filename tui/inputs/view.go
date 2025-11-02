package inputs

import (
	"strings"

	"github.com/razshare/frizzante/tui/config"
)

func (model *Model) View() string {
	var builder strings.Builder
	builder.WriteString(config.Styles.Menu.Render("│"))
	builder.WriteString(config.Styles.Menu.Render("⏣ " + model.Prompt))
	builder.WriteString("\n")
	builder.WriteString(config.Styles.Menu.Render("│"))
	builder.WriteString(config.Styles.Title.Render(model.TextInput.View()))
	builder.WriteString("\n")
	builder.WriteString(config.Styles.Menu.Render("│"))
	builder.WriteString(config.Styles.UserGuide.Render("enter submit"))
	if model.TextInput.View() != "" {
		builder.WriteString(config.Styles.UserGuide.Render(" • esc clear"))
	} else {
		builder.WriteString(config.Styles.UserGuide.Render(" • esc back"))
	}
	builder.WriteString("\n")
	return builder.String()
}
