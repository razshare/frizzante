package input

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/tui/config"
	"strings"
)

func (model *Model) Init() tea.Cmd {
	return textinput.Blink
}

func (model *Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch assert := message.(type) {
	case tea.KeyMsg:
		if assert.Type == tea.KeyCtrlC {
			return model, tea.Interrupt
		}

		if assert.Type == tea.KeyEsc {
			model.TextInput.Reset()

			if model.TextInput.Value() != "" {
				return model, nil
			}

			return model, tea.Quit
		}

		if assert.Type == tea.KeyEnter {
			return model, tea.Quit
		}
	}

	model.TextInput, cmd = model.TextInput.Update(message)
	return model, cmd
}

func (model *Model) View() string {
	var builder strings.Builder
	builder.WriteString(config.Styles.Menu.Render("│"))
	builder.WriteString(config.Styles.Menu.Render(" ⏣ " + model.Prompt))
	builder.WriteString("\n")
	builder.WriteString(config.Styles.Menu.Render("│"))
	builder.WriteString(config.Styles.Title.Render(" " + model.TextInput.View()))
	builder.WriteString("\n")
	builder.WriteString(config.Styles.UserGuide.Render("enter submit"))
	if model.TextInput.View() != "" {
		builder.WriteString(config.Styles.UserGuide.Render(" • esc clear"))
	} else {
		builder.WriteString(config.Styles.UserGuide.Render(" • esc back"))
	}
	builder.WriteString("\n")
	return builder.String()
}
