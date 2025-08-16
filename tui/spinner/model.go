package spinner

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func (m *Model) Init() tea.Cmd {
	return m.Spinner.Tick
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch k := msg.(type) {
	case tea.KeyMsg:
		if k.Type == tea.KeyCtrlC {
			os.Exit(0)
		}
	}

	var cmd tea.Cmd
	m.Spinner, cmd = m.Spinner.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	return fmt.Sprintf("%s %s\n", m.Spinner.View(), m.Message)
}
