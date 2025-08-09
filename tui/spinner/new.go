package spinner

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/tui/config"
)

func New(message string) *Spinner {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = config.Styles.Spinner

	m := Model{
		Spinner: s,
		Message: message,
	}

	return &Spinner{
		Model:   m,
		Program: tea.NewProgram(m),
	}
}
