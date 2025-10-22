package npmselect

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

		if assert.Type == tea.KeyEsc {
			model.Quitting = true
			model.Selected = make([]string, 0)
			return model, tea.Quit
		}

		if assert.Type == tea.KeyEnter {
			if len(model.Selected) == 0 && len(model.Search.Filtered) > 0 {
				val := model.Search.Filtered[model.Viewport.Cursor].Id
				model.Selected = append(model.Selected, val)
			}
			model.Confirmed = true
			return model, tea.Quit
		}

		if assert.Type == tea.KeySpace {
			if len(model.Search.Filtered) > 0 {
				val := model.Search.Filtered[model.Viewport.Cursor].Id
				if slices.Contains(model.Selected, val) {
					if i := slices.Index(model.Selected, val); i >= 0 {
						model.Selected = append(model.Selected[:i], model.Selected[i+1:]...)
					}
				} else {
					model.Selected = append(model.Selected, val)
				}
			}
			return model, nil
		}

		if assert.Type == tea.KeyUp || assert.Type == tea.KeyCtrlP {
			navigate.Apply(model.Search, model.Viewport, -1)
			return model, nil
		}

		if assert.Type == tea.KeyDown || assert.Type == tea.KeyCtrlN || assert.Type == tea.KeyTab {
			navigate.Apply(model.Search, model.Viewport, 1)
			return model, nil
		}

		if assert.Type == tea.KeyPgDown {
			navigate.Apply(model.Search, model.Viewport, model.Viewport.Visible)
			return model, nil
		}

		if assert.Type == tea.KeyPgUp {
			navigate.Apply(model.Search, model.Viewport, -model.Viewport.Visible)
			return model, nil
		}

		// Handle search input
		if len(assert.String()) == 1 || assert.Type == tea.KeyBackspace || assert.Type == tea.KeyCtrlH {
			if !model.Search.Active {
				model.Search.Active = true
				model.Search.Input.Focus()
			}

			var cmd tea.Cmd
			model.Debouncer.Reset(model.Debounce)
			model.Search.Input, cmd = model.Search.Input.Update(assert)
			model.LastQuery = model.Search.Input.Value()

			return model, tea.Batch(cmd, func() tea.Msg {
				<-model.Debouncer.C
				return DebouncedSearchMsg{Query: model.Search.Input.Value()}
			})
		}

	case DebouncedSearchMsg:
		if assert.Query != "" && assert.Query == model.Search.Input.Value() {
			model.Loading = true
			return model, Search(assert.Query)
		}

	case SearchResultMsg:
		model.Loading = false
		model.Error = assert.Error
		model.Packages = assert.Packages
		choices := make([]search.Choice, len(assert.Packages))
		for i, pkg := range assert.Packages {
			id := pkg.Name
			if pkg.Version != "" {
				id = fmt.Sprintf("%s@%s", pkg.Name, pkg.Version)
			}
			description := pkg.Description
			if len(description) > 50 {
				description = description[:50-3] + "..."
			}
			choices[i] = search.Choice{
				Id:          id,
				Description: description,
			}
		}
		model.Search.Choices = choices
		model.Search.Filtered = choices
		model.Viewport.Cursor = 0
		return model, nil
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

	if model.Loading {
		builder.WriteString("\n")
		builder.WriteString(config.Styles.UserGuide.Render("ⓘ  loading..."))
		builder.WriteString("\n")
		if model.Search.Active {
			builder.WriteString(config.Styles.UserGuide.Render(" • esc clear"))
		}
		return builder.String()
	} else if model.Error != nil {
		builder.WriteString("\n")
		builder.WriteString(config.Styles.Status(config.Colors.Error).Render("✗  " + model.Error.Error()))
		builder.WriteString("\n")
		if model.Search.Active {
			builder.WriteString(config.Styles.UserGuide.Render(" • esc clear"))
		}
		return builder.String()
	}

	filtered := len(model.Search.Filtered)
	if filtered == 0 {
		builder.WriteString("\n")
		builder.WriteString(config.Styles.Menu.Render("│"))
		builder.WriteString(config.Styles.UserGuide.Render("ⓘ  no matches found"))

		builder.WriteString("\n")

		builder.WriteString(config.Styles.Menu.Render("│"))
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
		builder.WriteString(config.Styles.Menu.Render("│"))
		builder.WriteString(config.Styles.Status(config.Colors.Muted).Render("↑ more above"))
		builder.WriteString("\n")
	}

	for i := model.Viewport.Offset; i < height; i++ {
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

	builder.WriteString(config.Styles.Menu.Render("│"))
	builder.WriteString(config.Styles.UserGuide.Render("↑ up • ↓ down • space select • enter continue"))

	if model.Search.Active {
		builder.WriteString(config.Styles.UserGuide.Render(" • esc clear"))
	} else {
		builder.WriteString(config.Styles.UserGuide.Render(" • esc back"))
	}

	return builder.String()
}
