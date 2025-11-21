package select_npm_packages

import (
	"fmt"
	"slices"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/tui/navigate"
	"github.com/razshare/frizzante/tui/search"
)

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

		if !model.Search.Active {
			model.Search.Active = true
		}

		model.Debouncer.Reset(model.Debounce)
		search.Apply(model.Search, model.Viewport, assert)
		model.LastQuery = model.Search.Value

		var cmd tea.Cmd
		return model, tea.Batch(cmd, func() tea.Msg {
			<-model.Debouncer.C
			return DebouncedSearchMsg{Query: model.Search.Value}
		})

	case DebouncedSearchMsg:
		if assert.Query != "" && assert.Query == model.Search.Value {
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
