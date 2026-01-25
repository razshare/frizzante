package hexviewer

import (
	"fmt"
	"strings"

	"github.com/razshare/frizzante/tui/configs"
)

func (model *Model) View() string {
	var builder strings.Builder
	builder.Grow(1024)
	builder.WriteString(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
	builder.WriteString(configs.Styles.Menu.Render(model.Prompt))
	if model.Search.Value != "" {
		builder.WriteString(configs.Styles.UserInput.Render(" ⁋/" + model.Search.Value))
	} else {
		builder.WriteString(configs.Styles.UserGuide.Render(" ⁋/type to search"))
	}
	filtered := len(model.Search.Filtered)
	if filtered == 0 {
		builder.WriteString("\n")
		builder.WriteString(configs.Styles.Menu.PaddingRight(2).Render("│"))
		builder.WriteString(configs.Styles.UserGuide.Render("ⓘ  no matches found"))
		builder.WriteString("\n")
		builder.WriteString(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
		builder.WriteString(configs.Styles.UserGuide.Render("↑ up • ↓ down"))
		if model.Search.Active {
			builder.WriteString(configs.Styles.UserGuide.Render(" • esc clear"))
		} else {
			builder.WriteString(configs.Styles.UserGuide.Render(" • esc back"))
		}
		return builder.String()
	}
	if filtered > 0 {
		builder.WriteString(configs.Styles.UserInput.Render(fmt.Sprintf(" (%d)", filtered)))
	}
	builder.WriteString("\n")
	height := model.Viewport.Offset + model.Viewport.Visible
	if height > filtered {
		height = filtered
	}
	if model.Viewport.Offset > 0 {
		builder.WriteString(configs.Styles.Menu.PaddingRight(2).Render("│"))
		builder.WriteString(configs.Styles.Status(configs.Colors.Muted).Render("↑ more above"))
		builder.WriteString("\n")
	}
	maxKeyWidth := len(fmt.Sprintf("%d.", height-1))
	for index := model.Viewport.Offset; index < height; index++ {
		key := fmt.Sprintf("%d.", index)
		padding := maxKeyWidth - len(key) + 3
		content := model.Search.Filtered[index].Id
		builder.WriteString(configs.Styles.Menu.PaddingRight(2).Render("│"))
		builder.WriteString(configs.Styles.Menu.PaddingRight(padding).Render(key))
		if model.Viewport.Cursor == index {
			builder.WriteString(configs.Styles.Selected.Render(content))
		} else if !strings.HasPrefix(strings.TrimSpace(content), "0x") {
			builder.WriteString(configs.Styles.Menu.Render(content))
		} else {
			builder.WriteString(configs.Styles.Item.Render(content))
		}
		builder.WriteString("\n")
	}
	if height < filtered {
		builder.WriteString(configs.Styles.Menu.PaddingRight(2).Render("│"))
		builder.WriteString(configs.Styles.Status(configs.Colors.Muted).Render("↓ more below"))
		builder.WriteString("\n")
	}
	builder.WriteString(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
	builder.WriteString(configs.Styles.UserGuide.Render("↑ up • ↓ down"))
	if model.Search.Active {
		builder.WriteString(configs.Styles.UserGuide.Render(" • esc clear"))
	} else {
		builder.WriteString(configs.Styles.UserGuide.Render(" • esc back"))
	}
	return builder.String()
}
