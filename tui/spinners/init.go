package spinners

import tea "github.com/charmbracelet/bubbletea"

func (model *Model) Init() tea.Cmd {
	return model.Spinner.Tick
}
