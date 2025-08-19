package multiselect

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/razshare/frizzante/tui/program"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/viewport"
)

func Send(chs []search.Choice, msg string) ([]string, error) {
	// Initialize the search input
	input := textinput.New()
	input.Width = 80
	model, err := program.Run(&Model{
		Search: &search.Search{
			Active:   false,
			Choices:  chs,
			Filtered: chs,
			Input:    input,
		},
		Viewport: &viewport.Viewport{
			Visible: 6,
			Start:   0,
			Cursor:  0,
		},
		Selected: make([]string, 0),
		Prompt:   msg,
	})

	if err != nil {
		return nil, err
	}

	return model.Selected, nil
}

func Sendf(opts []string, format string, vars ...any) []string {
	return Sendf(opts, fmt.Sprintf(format, vars...))
}
