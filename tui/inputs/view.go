package inputs

import (
	"strings"

	"github.com/razshare/frizzante/v2/tui/configs"
)

func (model *Model) View() string {
	var builder strings.Builder
	builder.WriteString(configs.Styles.Menu.Render("│"))
	builder.WriteString(configs.Styles.Menu.PaddingLeft(1).Render("⏣ " + model.Prompt))
	builder.WriteString("\n")
	builder.WriteString(configs.Styles.Menu.Render("│"))
	builder.WriteString(configs.Styles.UserGuide.PaddingLeft(3).Render("⁋/> "))
	if model.Value == "" && model.Prefix == "" {
		builder.WriteString(configs.Styles.UserGuide.Render("type here"))
	} else {
		if model.Prefix != "" {
			builder.WriteString(configs.Styles.UserGuide.Render(model.Prefix))
		}
		builder.WriteString(configs.Styles.UserInput.Render(model.Value))
	}
	builder.WriteString("\n")

	if model.ClipboardError != nil {
		builder.WriteString(configs.Styles.Menu.Render("│"))
		builder.WriteString(configs.Styles.ErrorLabel.Render("  ERROR  "))
		builder.WriteString(configs.Styles.ErrorText.Render(model.ClipboardError.Error()))
		builder.WriteString("\n")
	}

	builder.WriteString(configs.Styles.Menu.Render("│"))
	builder.WriteString(configs.Styles.UserGuide.Render("enter submit"))
	if model.Value != "" {
		builder.WriteString(configs.Styles.UserGuide.Render(" • esc clear"))
	} else {
		builder.WriteString(configs.Styles.UserGuide.Render(" • esc back"))
	}
	return builder.String()
}
