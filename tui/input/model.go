package input

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"main/config"
)

func (model Model) Init() tea.Cmd {
	return textinput.Blink
}

func (model Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch messageLocal := msg.(type) {
	case tea.KeyMsg:
		switch messageLocal.Type {
		case tea.KeyEnter:
			return model, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			model.Cancelled = true
			return model, tea.Quit
		}
	}

	model.TextInput, cmd = model.TextInput.Update(msg)
	return model, cmd
}

func (model Model) View() string {
	return fmt.Sprintf("\n%s\n\n%s\n\n%s",
		config.Styles.Title.Render(model.Prompt),
		model.TextInput.View(),
		lipgloss.NewStyle().Foreground(lipgloss.Color(config.Colors.Muted)).Render("(esc to quit)"))
}
