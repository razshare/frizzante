package select_many

import (
	"fmt"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/tui/config"
	"github.com/razshare/frizzante/tui/navigate"
	"github.com/razshare/frizzante/tui/search"
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
				if !slices.Contains(model.Selected, value) {
					model.Selected = append(model.Selected, value)
				}
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
				return model, nil
			}

			model.Selected = make([]string, 0)
			return model, tea.Quit
		}

		if assert.Type == tea.KeyUp || assert.Type == tea.KeyCtrlP || assert.Type == tea.KeyShiftTab || assert.Type == tea.KeyCtrlPgUp {
			navigate.Apply(model.Search, model.Viewport, -1)
			return model, nil
		}

		if assert.Type == tea.KeyDown || assert.Type == tea.KeyCtrlN || assert.Type == tea.KeyTab || assert.Type == tea.KeyCtrlPgDown {
			navigate.Apply(model.Search, model.Viewport, 1)
			return model, nil
		}

		if assert.Type == tea.KeyPgUp {
			navigate.Apply(model.Search, model.Viewport, -model.Viewport.Visible)
			return model, nil
		}

		if assert.Type == tea.KeyPgDown {
			navigate.Apply(model.Search, model.Viewport, model.Viewport.Visible)
			return model, nil
		}

		if assert.Type == tea.KeyHome {
			home := model.Viewport.Offset
			if model.Viewport.Cursor == home {
				navigate.Apply(model.Search, model.Viewport, -model.Viewport.Visible)
			} else {
				model.Viewport.Cursor = home
				navigate.Apply(model.Search, model.Viewport, 0)
			}
			return model, nil
		}

		if assert.Type == tea.KeyEnd {
			end := model.Viewport.Offset + model.Viewport.Visible - 1
			if model.Viewport.Cursor == end {
				navigate.Apply(model.Search, model.Viewport, model.Viewport.Visible)
			} else {
				model.Viewport.Cursor = end
				navigate.Apply(model.Search, model.Viewport, 0)
			}
			return model, nil
		}

		// vscode has a weird bug where it will send a "ctrl+w" whenever the user presses "backspace" in the integrated terminal,
		// so we're including tea.KeyCtrlW to try to fix that for the user.
		// https://stackoverflow.com/questions/52806758/visual-studio-code-ctrlbackspace-not-working-in-integrated-terminal
		if assert.Type == tea.KeyBackspace || assert.Type == tea.KeyCtrlH || assert.Type == tea.KeyCtrlW {
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

		builder.WriteString(config.Styles.Menu.PaddingRight(1).Render("│"))
		builder.WriteString(config.Styles.UserGuide.Render("↑ up • ↓ down • space select • enter continue"))

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
		builder.WriteString(config.Styles.Menu.PaddingRight(2).Render("│"))
		builder.WriteString(config.Styles.Status(config.Colors.Muted).Render("↓ more below"))
		builder.WriteString("\n")
	}

	builder.WriteString(config.Styles.Menu.PaddingRight(1).Render("│"))
	builder.WriteString(config.Styles.UserGuide.Render("↑ up • ↓ down • space select • enter continue"))

	if model.Search.Active {
		builder.WriteString(config.Styles.UserGuide.Render(" • esc clear"))
	} else {
		builder.WriteString(config.Styles.UserGuide.Render(" • esc back"))
	}

	return builder.String()
}
