package spinner

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/tui/config"
)

var _ = spinner.Spinner{
	Frames: []string{"䷀", "䷫", "䷠", "䷋", "䷓", "䷚", "䷨", "䷙", "䷍", "䷍", "䷍", "䷡", "䷪"},
	FPS:    time.Second / 6,
}

var _star = spinner.Spinner{
	Frames: []string{"⯌", "⯍", "⯎", "⯏"},
	FPS:    time.Second / 4,
}

func New(message string) *Spinner {
	spin := spinner.New()
	spin.Spinner = _star
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
