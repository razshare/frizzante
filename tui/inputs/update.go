package inputs

import tea "github.com/charmbracelet/bubbletea"

func (model *Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch assert := message.(type) {
	case tea.KeyMsg:
		if assert.Type == tea.KeyCtrlC {
			return model, tea.Interrupt
		}

		if assert.Type == tea.KeyEsc {
			if model.TextInput.Value() != "" {
				model.TextInput.Reset()
				return model, nil
			}

			return model, tea.Quit
		}

		if assert.Type == tea.KeyEnter {
			return model, tea.Quit
		}
	}

	model.TextInput, cmd = model.TextInput.Update(message)
	return model, cmd
}
