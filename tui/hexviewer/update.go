package hexviewer

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
		if assert.Type == tea.KeyEsc {
			if model.Search.Active {
				search.Reset(model.Search, model.Viewport)
				return model, nil
			}
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
			}
			return model, search.Apply(model.Search, model.Viewport, assert)
		}
	}
	return model, nil
}
