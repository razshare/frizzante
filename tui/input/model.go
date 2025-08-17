package input

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/tui/config"
	"strings"
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

		if k.Type == tea.KeyEsc {
			m.TextInput.Reset()
			return m, nil
		}

		if k.Type == tea.KeyEnter {
			return m, tea.Quit
		}
	}

	m.TextInput, cmd = m.TextInput.Update(msg)
	return m, cmd
}

func (m *Model) View() string {
	var sb strings.Builder
	sb.WriteString(config.Styles.Menu.Render("│"))
	sb.WriteString(config.Styles.Menu.Render(" ⏣ " + m.Prompt))
	sb.WriteString("\n")
	sb.WriteString(config.Styles.Menu.Render("│"))
	sb.WriteString(config.Styles.Title.Render(" " + m.TextInput.View()))
	sb.WriteString("\n")
	sb.WriteString(config.Styles.UserGuide.Render("enter submit • esc clear"))
	sb.WriteString("\n")
	return sb.String()
}
