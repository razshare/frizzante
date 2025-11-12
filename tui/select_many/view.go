package select_many

import (
	"fmt"
	"slices"
	"strings"

	"github.com/razshare/frizzante/tui/configs"
)

func (model *Model) View() string {
	var builder strings.Builder
	builder.Grow(1024)

	builder.WriteString(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
	builder.WriteString(configs.Styles.Menu.Render(model.Prompt))

	builder.WriteString(configs.Styles.UserGuide.PaddingLeft(1).Render("⁋/> "))
	if model.Search.Input.Value() != "" {
		builder.WriteString(configs.Styles.UserInput.Render(model.Search.Input.Value()))
	} else {
		builder.WriteString(configs.Styles.UserGuide.Render("type to search"))
	}

	filtered := len(model.Search.Filtered)
	if filtered == 0 {
		builder.WriteString("\n")
		builder.WriteString(configs.Styles.Menu.PaddingRight(2).Render("│"))
		builder.WriteString(configs.Styles.UserGuide.Render("ⓘ  no matches found"))

		builder.WriteString("\n")

		builder.WriteString(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
		builder.WriteString(configs.Styles.UserGuide.Render("↑ up • ↓ down • space select • enter continue"))

		if model.Search.Active {
			builder.WriteString(configs.Styles.UserGuide.Render(" • esc clear"))
		} else {
			builder.WriteString(configs.Styles.UserGuide.Render(" • esc back"))
		}

		return builder.String()
	} else {
		if filtered > 0 {
			builder.WriteString(configs.Styles.UserInput.Render(fmt.Sprintf(" (%d)", filtered)))
		}
		builder.WriteString("\n")
	}

	height := model.Viewport.Offset + model.Viewport.Visible
	if height > filtered {
		height = filtered
	}

	if model.Viewport.Offset > 0 {
		builder.WriteString(configs.Styles.Menu.PaddingRight(2).Render("│"))
		builder.WriteString(configs.Styles.UserGuide.Render("↑ more above"))
		builder.WriteString("\n")
	}

	for i := model.Viewport.Offset; i < height; i++ {
		builder.WriteString(configs.Styles.Menu.PaddingRight(2).Render("│"))
		if model.Viewport.Cursor == i {
			if slices.Contains(model.Selected, model.Search.Filtered[i].Id) {
				builder.WriteString(configs.Styles.Selected.Render("● " + model.Search.Filtered[i].Id))
			} else {
				builder.WriteString(configs.Styles.Selected.Render("◉ " + model.Search.Filtered[i].Id))
			}

			j := slices.Index(model.Search.Choices, model.Search.Filtered[i])
			if j >= 0 && model.Search.Choices[j].Description != "" {
				builder.WriteString(configs.Styles.UserGuide.Render("  ⇢  " + model.Search.Choices[j].Description))
			}
		} else if slices.Contains(model.Selected, model.Search.Filtered[i].Id) {
			builder.WriteString(configs.Styles.Item.Render("● " + model.Search.Filtered[i].Id))
		} else {
			builder.WriteString(configs.Styles.Item.Render("○ " + model.Search.Filtered[i].Id))
		}
		builder.WriteString("\n")
	}

	if height < filtered {
		builder.WriteString(configs.Styles.Menu.PaddingRight(2).Render("│"))
		builder.WriteString(configs.Styles.UserGuide.Render("↓ more below"))
		builder.WriteString("\n")
	}

	builder.WriteString(configs.Styles.Menu.PaddingRight(1).Render("⎚"))
	builder.WriteString(configs.Styles.UserGuide.Render("↑ up • ↓ down • space select • enter continue"))

	if model.Search.Active {
		builder.WriteString(configs.Styles.UserGuide.Render(" • esc clear"))
	} else {
		builder.WriteString(configs.Styles.UserGuide.Render(" • esc back"))
	}

	return builder.String()
}
