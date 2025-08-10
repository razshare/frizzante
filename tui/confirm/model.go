package confirm

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/tui/config"
	"strings"
)

func (model Model) Init() tea.Cmd {
	return nil
}

func (model Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch k := msg.(type) {
	case tea.KeyMsg:
		if k.Type == tea.KeyCtrlC {
			model.Confirmed = false
			return model, tea.Quit
		}

		if k.Type == tea.KeyCtrlC {
			model.Confirmed = model.DefaultValue
			return model, tea.Quit
		}

		if strings.ToLower(k.String()) == "y" {
			model.Confirmed = true
			return model, tea.Quit
		}

		if strings.ToLower(k.String()) == "n" {
			model.Confirmed = false
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
