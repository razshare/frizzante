package input

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/tui/config"
)

func (model Model) Init() tea.Cmd {
	return textinput.Blink
}

func (model Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch k := msg.(type) {
	case tea.KeyMsg:
		if k.Type == tea.KeyEnter {
			return model, tea.Quit
		}

		if k.Type == tea.KeyCtrlC || k.Type == tea.KeyEsc {
			model.Cancelled = true
			return model, tea.Quit
		}
	}

	model.TextInput, cmd = model.TextInput.Update(msg)
	return model, cmd
}

func (model Model) View() string {
	return fmt.Sprintf(
		"\n%s\n\n%s\n\n%s",
		config.Styles.Title.Render(model.Prompt),
		model.TextInput.View(),
		config.Styles.Suggestion.Render("(esc to quit)"),
	)
}
