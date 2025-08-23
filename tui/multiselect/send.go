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
		Prompt:   msg,
		Viewport: &viewport.Viewport{Visible: 6},
		Selected: make([]string, 0),
		Search: &search.Search{
			Choices:  chs,
			Filtered: chs,
			Input:    input,
		},
	})

	if err != nil {
		return nil, err
	}

	return model.Selected, nil
}

func Sendf(chs []search.Choice, format string, vars ...any) ([]string, error) {
	return Send(chs, fmt.Sprintf(format, vars...))
}
