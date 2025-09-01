package spinner

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Spinner       spinner.Model
	Message       string
	QuitMessage   string
	SoftInterrupt bool
}

type Spinner struct {
	Model   *Model
	Program *tea.Program
	Done    chan bool
}
