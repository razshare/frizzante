package select_one

import (
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
		return model, search.Apply(model.Search, model.Viewport, assert)
	}
	return model, nil
}
