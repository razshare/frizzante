package inputs

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (model *Model) Init() tea.Cmd {
	return textinput.Blink
}
