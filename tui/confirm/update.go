package confirm

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (model *Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch assert := message.(type) {
	case tea.KeyMsg:
		if assert.Type == tea.KeyEsc {
			model.Confirmed = false
			return model, tea.Quit
		}

		if assert.Type == tea.KeyCtrlC {
			model.Confirmed = false
			return model, tea.Interrupt
		}

		if strings.ToLower(assert.String()) == "y" {
			model.Confirmed = true
			return model, tea.Quit
		}

		if strings.ToLower(assert.String()) == "n" {
			model.Confirmed = false
			return model, tea.Quit
		}

		if assert.Type == tea.KeyLeft ||
			assert.Type == tea.KeyCtrlLeft ||
			assert.Type == tea.KeyRight ||
			assert.Type == tea.KeyCtrlRight ||
			assert.Type == tea.KeyTab {
			model.Confirmed = !model.Confirmed
			return model, nil
		}

		if assert.Type == tea.KeyEnter {
			return model, tea.Quit
		}
	}
	return model, nil
}
