package select_one

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/razshare/frizzante/tui/program"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/viewport"
)

func Send(choices []search.Choice, message string) (selected string, err error) {
	input := textinput.New()
	input.Width = 80
	var model *Model
	if model, err = program.Run(&Model{
		Prompt:   message,
		Viewport: &viewport.Viewport{Visible: 6},
		Search: &search.Search{
			Choices:  choices,
			Filtered: choices,
		},
	}); err != nil {
		return
	}
	selected = model.Selected
	return
}
