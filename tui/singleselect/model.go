package singleselect

import (
	"fmt"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/tui/config"
	"github.com/razshare/frizzante/tui/navigate"
	"github.com/razshare/frizzante/tui/search"
)

func (model *Model) Init() tea.Cmd {
	return textinput.Blink
}

func (model *Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch assert := message.(type) {
	case tea.KeyMsg:
		if assert.Type == tea.KeyCtrlC {
			return model, tea.Interrupt
		}

		if assert.Type == tea.KeyEnter {
			if len(model.Search.Filtered) > 0 {
				model.Selected = model.Search.Filtered[model.Viewport.Cursor].Id
				return model, tea.Quit
			}
		}

		if assert.Type == tea.KeyEsc {
			if model.Search.Active {
				search.Reset(model.Search, model.Viewport)
				return model, nil
			}

			model.Selected = ""
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
		builder.WriteString(config.Styles.Menu.Render("│"))
		builder.WriteString(config.Styles.UserGuide.Render("ⓘ  no matches found"))

		builder.WriteString("\n")

		builder.WriteString(config.Styles.Menu.Render("│"))
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
		builder.WriteString(config.Styles.Menu.Render("│"))
		builder.WriteString(config.Styles.Status(config.Colors.Muted).Render("↑ more above"))
		builder.WriteString("\n")
	}

	for i := model.Viewport.Offset; i < height; i++ {
		builder.WriteString(config.Styles.Menu.Render("│"))
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
		builder.WriteString(config.Styles.Menu.Render("│"))
		builder.WriteString(config.Styles.Status(config.Colors.Muted).Render("↓ more below"))
		builder.WriteString("\n")
	}

	builder.WriteString(config.Styles.Menu.Render("│"))
	builder.WriteString(config.Styles.UserGuide.Render("↑ up • ↓ down • enter submit"))

	if model.Search.Active {
		builder.WriteString(config.Styles.UserGuide.Render(" • esc clear"))
	} else {
		builder.WriteString(config.Styles.UserGuide.Render(" • esc back"))
	}

	return builder.String()
}
