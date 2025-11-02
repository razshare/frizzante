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
			return model, tea.Interrupt
		}

		if assert.Type == tea.KeyEnter {
			model.Confirmed = model.DefaultValue
			return model, tea.Quit
		}

		if strings.ToLower(assert.String()) == "y" {
			model.Confirmed = true
			return model, tea.Quit
		}

		if strings.ToLower(assert.String()) == "n" {
			model.Confirmed = false
			return model, tea.Quit
		}
	}
	return model, nil
}
