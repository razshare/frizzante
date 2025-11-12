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
			if model.Value != "" {
				model.Value = ""
				return model, nil
			}

			return model, tea.Quit
		}

		if assert.Type == tea.KeyEnter {
			return model, tea.Quit
		}

		// vscode has a weird bug where it will send a "ctrl+w" whenever the user presses "backspace" in the integrated terminal,
		// so we're including tea.KeyCtrlW to try to fix that for the user.
		// https://stackoverflow.com/questions/52806758/visual-studio-code-ctrlbackspace-not-working-in-integrated-terminal
		if assert.Type == tea.KeyBackspace || assert.Type == tea.KeyCtrlH || assert.Type == tea.KeyCtrlW {
			length := len(model.Value)
			if length > 0 {
				model.Value = model.Value[:length-1]
			}
			return model, nil
		}

		model.Value += assert.String()
	}

	return model, cmd
}
