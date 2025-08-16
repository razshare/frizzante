package confirm

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/tui/config"
	"strings"
)

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch k := msg.(type) {
	case tea.KeyMsg:
		if k.Type == tea.KeyCtrlC {
			if m.SoftInterrupt {
				return m, tea.Quit
			}
			return m, tea.Interrupt
		}

		if k.Type == tea.KeyEnter {
			m.Confirmed = m.DefaultValue
			return m, tea.Quit
		}

		if strings.ToLower(k.String()) == "y" {
			m.Confirmed = true
			return m, tea.Quit
		}

		if strings.ToLower(k.String()) == "n" {
			m.Confirmed = false
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *Model) View() string {
	if m.DefaultValue {
		return config.Styles.Title.Render(m.Prompt, "(Y/n)")
	}
	return config.Styles.Title.Render(m.Prompt, "(y/N)")
}
