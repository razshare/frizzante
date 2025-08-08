package spinner

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

func (model Model) Init() tea.Cmd {
	return model.Spinner.Tick
}

func (model Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	model.Spinner, cmd = model.Spinner.Update(msg)
	return model, cmd
}

func (model Model) View() string {
	return fmt.Sprintf("%s %s", model.Spinner.View(), model.Message)
}
