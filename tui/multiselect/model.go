package multiselect

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/tui/config"
	"github.com/razshare/frizzante/tui/navigate"
	"github.com/razshare/frizzante/tui/search"
	"slices"
	"strings"
)

func (model *Model) Init() tea.Cmd {
	return nil
}

func (model *Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch assert := message.(type) {
	case tea.KeyMsg:
		if assert.Type == tea.KeyCtrlC {
			return model, tea.Interrupt
		}

		if assert.Type == tea.KeyEnter {
			if len(model.Selected) == 0 && len(model.Search.Filtered) > 0 {
				value := model.Search.Filtered[model.Viewport.Cursor].Id
				model.Selected = append(model.Selected, value)
			}
			return model, tea.Quit
		}

		if assert.Type == tea.KeySpace {
			if len(model.Search.Filtered) > 0 {
				value := model.Search.Filtered[model.Viewport.Cursor].Id
				if slices.Contains(model.Selected, value) {
					if i := slices.Index(model.Selected, value); i >= 0 {
						model.Selected = append(model.Selected[:i], model.Selected[i+1:]...)
					}
				} else {
					model.Selected = append(model.Selected, value)
				}
			}
			return model, nil
		}

		if assert.Type == tea.KeyEsc {
			if model.Search.Active {
				search.Reset(model.Search, model.Viewport)
			}

			model.Selected = make([]string, 0)
			return model, tea.Quit
		}

		if assert.Type == tea.KeyUp || assert.Type == tea.KeyCtrlP {
			navigate.Apply(model.Search, model.Viewport, -1)
			return model, nil
		}

		if assert.Type == tea.KeyDown || assert.Type == tea.KeyCtrlN || assert.Type == tea.KeyTab {
			navigate.Apply(model.Search, model.Viewport, 1)
			return model, nil
		}

		if assert.Type == tea.KeyBackspace || assert.Type == tea.KeyCtrlH {
			if model.Search.Active {
				return model, search.Apply(model.Search, model.Viewport, assert)
			}
		}

		if len(assert.String()) == 1 {
			if !model.Search.Active {
				model.Search.Active = true
				model.Search.Input.Focus()
			}
			return model, search.Apply(model.Search, model.Viewport, assert)
		}
	}
	return model, nil
}

func (model *Model) View() string {
	var builder strings.Builder
	builder.Grow(1024)

	builder.WriteString(config.Styles.Menu.Render(model.Prompt))

	if model.Search.Input.Value() != "" {
		builder.WriteString(config.Styles.UserInput.Render(" ⁋/" + model.Search.Input.Value()))
	} else {
		builder.WriteString(config.Styles.UserGuide.Render(" ⁋/type to search"))
	}

	builder.WriteString("\n")

	filtered := len(model.Search.Filtered)
	if filtered == 0 {
		builder.WriteString(config.Styles.Menu.Render("│"))
		builder.WriteString(config.Styles.UserGuide.PaddingLeft(1).Render("ⓘ  no matches found"))

		builder.WriteString("\n")

		builder.WriteString(config.Styles.UserGuide.Render("↑ up • ↓ down • space select • enter continue"))

		if model.Search.Active {
			builder.WriteString(config.Styles.UserGuide.Render(" • esc clear"))
		} else {
			builder.WriteString(config.Styles.UserGuide.Render(" • esc back"))
		}

		return builder.String()
	}

	height := model.Viewport.Start + model.Viewport.Visible
	if height > filtered {
		height = filtered
	}

	if model.Viewport.Start > 0 {
		builder.WriteString(config.Styles.Menu.Render("│"))
		builder.WriteString(config.Styles.Status(config.Colors.Muted).Render("↑ more above"))
		builder.WriteString("\n")
	}

	for i := model.Viewport.Start; i < height; i++ {
		builder.WriteString(config.Styles.Menu.Render("│"))
		if model.Viewport.Cursor == i {
			if slices.Contains(model.Selected, model.Search.Filtered[i].Id) {
				builder.WriteString(config.Styles.Selected.Render("● " + model.Search.Filtered[i].Id))
			} else {
				builder.WriteString(config.Styles.Selected.Render("◉ " + model.Search.Filtered[i].Id))
			}

			j := slices.Index(model.Search.Choices, model.Search.Filtered[i])
			if j >= 0 && model.Search.Choices[j].Description != "" {
				builder.WriteString(config.Styles.UserGuide.Render("  ⇢  " + model.Search.Choices[j].Description))
			}
		} else if slices.Contains(model.Selected, model.Search.Filtered[i].Id) {
			builder.WriteString(config.Styles.Item.Render("● " + model.Search.Filtered[i].Id))
		} else {
			builder.WriteString(config.Styles.Item.Render("○ " + model.Search.Filtered[i].Id))
		}
		builder.WriteString("\n")
	}

	if height < filtered {
		builder.WriteString(config.Styles.Menu.Render("│"))
		builder.WriteString(config.Styles.Status(config.Colors.Muted).Render("↓ more below"))
		builder.WriteString("\n")
	}

	builder.WriteString(config.Styles.UserGuide.Render("↑ up • ↓ down • space select • enter continue"))

	if model.Search.Active {
		builder.WriteString(config.Styles.UserGuide.Render(" • esc clear"))
	} else {
		builder.WriteString(config.Styles.UserGuide.Render(" • esc back"))
	}

	return builder.String()
}
