package spinners

import tea "github.com/charmbracelet/bubbletea"

func (model *Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch assert := message.(type) {
	case tea.KeyMsg:
		if assert.Type == tea.KeyCtrlC {
			return model, tea.Interrupt
		}
	}
	var cmd tea.Cmd
	model.Spinner, cmd = model.Spinner.Update(message)
	return model, cmd
}
