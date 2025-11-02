package spinners

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/tui/config"
)

func New(message string) *Spinner {
	spin := spinner.New()
	spin.Spinner = Animation
	spin.Style = config.Styles.Menu

	model := &Model{
		Spinner: spin,
		Message: message,
	}

	return &Spinner{
		Model:   model,
		Program: tea.NewProgram(model),
	}
}

func Newf(format string, vars ...any) *Spinner {
	return New(fmt.Sprintf(format, vars...))
}
