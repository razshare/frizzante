package hexviewer

import (
	"fmt"
	"strings"

	"github.com/razshare/frizzante/tui/config"
)

func (model *Model) View() string {
	var builder strings.Builder
	builder.Grow(1024)

	builder.WriteString(config.Styles.Menu.PaddingRight(1).Render("⎚"))
	builder.WriteString(config.Styles.Menu.Render(model.Prompt))

	if model.Search.Input.Value() != "" {
		builder.WriteString(config.Styles.UserInput.Render(" ⁋/" + model.Search.Input.Value()))
	} else {
		builder.WriteString(config.Styles.UserGuide.Render(" ⁋/type to search"))
	}

	filtered := len(model.Search.Filtered)
	if filtered == 0 {
		builder.WriteString("\n")
		builder.WriteString(config.Styles.Menu.PaddingRight(2).Render("│"))
		builder.WriteString(config.Styles.UserGuide.Render("ⓘ  no matches found"))

		builder.WriteString("\n")

		builder.WriteString(config.Styles.Menu.PaddingRight(1).Render("⎚"))
		builder.WriteString(config.Styles.UserGuide.Render("↑ up • ↓ down"))

		if model.Search.Active {
			builder.WriteString(config.Styles.UserGuide.Render(" • esc clear"))
		} else {
			builder.WriteString(config.Styles.UserGuide.Render(" • esc back"))
		}

		return builder.String()
	} else {
		if filtered > 0 {
			builder.WriteString(config.Styles.UserInput.Render(fmt.Sprintf(" (%d)", filtered)))
		}
		builder.WriteString("\n")
	}

	height := model.Viewport.Offset + model.Viewport.Visible
	if height > filtered {
		height = filtered
	}

	if model.Viewport.Offset > 0 {
		builder.WriteString(config.Styles.Menu.PaddingRight(2).Render("│"))
		builder.WriteString(config.Styles.Status(config.Colors.Muted).Render("↑ more above"))
		builder.WriteString("\n")
	}

	maxKeyWidth := len(fmt.Sprintf("%d.", height-1))

	for index := model.Viewport.Offset; index < height; index++ {
		key := fmt.Sprintf("%d.", index)
		padding := maxKeyWidth - len(key) + 3
		content := model.Search.Filtered[index].Id

		builder.WriteString(config.Styles.Menu.PaddingRight(2).Render("│"))
		builder.WriteString(config.Styles.Menu.PaddingRight(padding).Render(key))

		if model.Viewport.Cursor == index {
			builder.WriteString(config.Styles.Selected.Render(content))
		} else if !strings.HasPrefix(strings.TrimSpace(content), "0x") {
			builder.WriteString(config.Styles.Menu.Render(content))
		} else {
			builder.WriteString(config.Styles.Item.Render(content))
		}

		builder.WriteString("\n")
	}

	if height < filtered {
		builder.WriteString(config.Styles.Menu.PaddingRight(2).Render("│"))
		builder.WriteString(config.Styles.Status(config.Colors.Muted).Render("↓ more below"))
		builder.WriteString("\n")
	}

	builder.WriteString(config.Styles.Menu.PaddingRight(1).Render("⎚"))
	builder.WriteString(config.Styles.UserGuide.Render("↑ up • ↓ down"))

	if model.Search.Active {
		builder.WriteString(config.Styles.UserGuide.Render(" • esc clear"))
	} else {
		builder.WriteString(config.Styles.UserGuide.Render(" • esc back"))
	}

	return builder.String()
}
