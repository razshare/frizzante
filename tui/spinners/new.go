package spinners

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/tui/configs"
)

func New(message string) *Spinner {
	spin := spinner.New()
	spin.Spinner = Animation
	spin.Style = configs.Styles.Menu
	model := &Model{
		Spinner: spin,
		Message: message,
	}
	return &Spinner{
		Model:   model,
		Program: tea.NewProgram(model),
	}
}
