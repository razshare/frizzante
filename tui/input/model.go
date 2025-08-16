package input

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/tui/config"
)

func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch k := msg.(type) {
	case tea.KeyMsg:
		if k.Type == tea.KeyCtrlC {
			if m.SoftInterrupt {
				return m, tea.Quit
			}
			return m, tea.Interrupt
		}

		if k.Type == tea.KeyEnter {
			return m, tea.Quit
		}
	}

	m.TextInput, cmd = m.TextInput.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	return fmt.Sprintf(
		"\n%s\n\n%s\n\n%s",
		config.Styles.Title.Render(m.Prompt),
		m.TextInput.View(),
		config.Styles.UserGuide.Render("[Esc] = clear"),
	)
}
