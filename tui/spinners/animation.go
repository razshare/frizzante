package spinners

import (
	"time"

	"github.com/charmbracelet/bubbles/spinner"
)

var Animation = spinner.Spinner{
	Frames: []string{"⯌ ", "⯏ ", "⯎ ", "⯍ "},
	FPS:    time.Second / 4,
}
