package spinners

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/tui/config"
)

//var _ = spinner.Spinner{
//	Frames: []string{"䷀", "䷫", "䷠", "䷋", "䷓", "䷚", "䷨", "䷙", "䷍", "䷍", "䷍", "䷡", "䷪"},
//	FPS:    time.Second / 6,
//}

var Star = spinner.Spinner{
	Frames: []string{"⯌ ", "⯏ ", "⯎ ", "⯍ "},
	FPS:    time.Second / 4,
}

func New(message string) *Spinner {
	spin := spinner.New()
	spin.Spinner = Star
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
