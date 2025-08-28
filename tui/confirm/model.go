package confirm

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/tui/config"
	"strings"
)

func (model *Model) Init() tea.Cmd {
	return nil
}

func (model *Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch assert := message.(type) {
	case tea.KeyMsg:
		if assert.Type == tea.KeyEsc {
			model.Confirmed = false
			return model, tea.Quit
		}

		if assert.Type == tea.KeyCtrlC {
			return model, tea.Interrupt
		}

		if assert.Type == tea.KeyEnter {
			model.Confirmed = model.DefaultValue
			return model, tea.Quit
		}

		if strings.ToLower(assert.String()) == "y" {
			model.Confirmed = true
			return model, tea.Quit
		}

		if strings.ToLower(assert.String()) == "n" {
			model.Confirmed = false
			return model, tea.Quit
		}
	}
	return model, nil
}

func (model *Model) View() string {
	if model.DefaultValue {
		return config.Styles.Title.Render(model.Prompt, "(Y/n)")
	}
	return config.Styles.Title.Render(model.Prompt, "(y/N)")
}
