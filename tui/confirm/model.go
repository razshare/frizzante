package confirm

import (
	tea "github.com/charmbracelet/bubbletea"
	"main/config"
)

func (model Model) Init() tea.Cmd {
	return nil
}

func (model Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch messageLocal := msg.(type) {
	case tea.KeyMsg:
		switch messageLocal.String() {
		case "y", "Y":
			model.Confirmed = true
			return model, tea.Quit
		case "n", "N":
			model.Confirmed = false
			return model, tea.Quit
		case "ctrl+c", "esc":
			model.Confirmed = false
			return model, tea.Quit
		case "enter":
			model.Confirmed = model.DefaultValue
			return model, tea.Quit
		}
	}
	return model, nil
}

func (model Model) View() string {
	defaultHint := "n"
	if model.DefaultValue {
		defaultHint = "Y"
	}
	return "\n" + config.Styles.Title.Render(model.Prompt) + "\n(Y/n) [default: " + defaultHint + "]"
}
