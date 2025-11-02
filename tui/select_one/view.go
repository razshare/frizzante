package select_one

import (
	"fmt"
	"slices"
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
		builder.WriteString(config.Styles.UserGuide.Render("↑ up • ↓ down • enter submit"))

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

	for i := model.Viewport.Offset; i < height; i++ {
		builder.WriteString(config.Styles.Menu.PaddingRight(2).Render("│"))
		if model.Viewport.Cursor == i {
			builder.WriteString(config.Styles.Selected.Render("● " + model.Search.Filtered[i].Id))
			j := slices.Index(model.Search.Choices, model.Search.Filtered[i])
			if j >= 0 && model.Search.Choices[j].Description != "" {
				builder.WriteString(config.Styles.UserGuide.Render("  ⇢  " + model.Search.Choices[j].Description))
			}
		} else {
			builder.WriteString(config.Styles.Item.Render("○ " + model.Search.Filtered[i].Id))
		}
		builder.WriteString("\n")
	}

	if height < filtered {
		builder.WriteString(config.Styles.Menu.PaddingRight(2).Render("│"))
		builder.WriteString(config.Styles.Status(config.Colors.Muted).Render("↓ more below"))
		builder.WriteString("\n")
	}
	builder.WriteString(config.Styles.Menu.PaddingRight(1).Render("⎚"))
	builder.WriteString(config.Styles.UserGuide.Render("↑ up • ↓ down • enter submit"))

	if model.Search.Active {
		builder.WriteString(config.Styles.UserGuide.Render(" • esc clear"))
	} else {
		builder.WriteString(config.Styles.UserGuide.Render(" • esc back"))
	}

	return builder.String()
}
