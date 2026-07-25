package select_many

import (
	"github.com/razshare/frizzante/v2/tui/program"
	"github.com/razshare/frizzante/v2/tui/search"
	"github.com/razshare/frizzante/v2/tui/viewport"
)

func Send(choices []search.Choice, message string) (selected []string, err error) {
	var model *Model
	if model, err = program.Run(&Model{
		Prompt:   message,
		Viewport: &viewport.Viewport{Visible: 6},
		Selected: make([]string, 0),
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
