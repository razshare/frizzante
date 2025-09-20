package spinner

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/tui/config"
)

var _wall = spinner.Spinner{
	Frames: []string{"䷀", "䷫", "䷠", "䷋", "䷓", "䷚", "䷨", "䷙", "䷍", "䷍", "䷍", "䷡", "䷪"},
	FPS:    time.Second / 6,
}

func New(message string) *Spinner {
	spin := spinner.New()
	spin.Spinner = _wall
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
