package npmselect

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/razshare/frizzante/tui/program"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/viewport"
)

func Send() (selected []string, err error) {
	input := textinput.New()
	input.Width = 80
	var model *Model
	if model, err = program.Run(&Model{
		Prompt:    "search npm packages",
		Viewport:  &viewport.Viewport{Visible: 6},
		Selected:  make([]string, 0),
		Debounce:  time.Second,
		Debouncer: time.NewTimer(time.Second),
		Search: &search.Search{
			Choices:  []search.Choice{},
			Filtered: []search.Choice{},
			Input:    input,
		},
	}); err != nil {
		return
	}

	selected = model.Selected

	return
}
