package select_many

import (
	"fmt"

	"github.com/razshare/frizzante/tui/program"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/viewport"
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

func Sendf(choices []search.Choice, format string, vars ...any) (selected []string, err error) {
	return Send(choices, fmt.Sprintf(format, vars...))
}
