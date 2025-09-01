package spinner

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/tui/config"
)

func New(message string) *Spinner {
	spin := spinner.New()
	spin.Spinner = spinner.Moon
	spin.Style = config.Styles.Spinner

	model := &Model{
		Spinner: spin,
		Message: message,
	}

	return &Spinner{
		Model:   model,
		Program: tea.NewProgram(model),
	}
}
