package spinner

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func (model *Model) Init() tea.Cmd {
	return model.Spinner.Tick
}

func (model *Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch assert := message.(type) {
	case tea.KeyMsg:
		if assert.Type == tea.KeyEsc {
			return model, tea.Quit
		}

		if assert.Type == tea.KeyCtrlC {
			if model.SoftInterrupt {
				return model, tea.Quit
			}
			return model, tea.Interrupt
		}
	}

	var cmd tea.Cmd
	model.Spinner, cmd = model.Spinner.Update(message)
	return model, cmd
}

func (model *Model) View() string {
	return fmt.Sprintf("%s %s\n", model.Spinner.View(), model.Message)
}
