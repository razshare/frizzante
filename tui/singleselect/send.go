package singleselect

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/tui/program"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/viewport"
)

func Send(chs []search.Choice, msg string) (string, error) {

	//c := len(opts)
	//chs := make([]string, c)
	//dsc := make([]string, c)

	//for i, option := range opts {
	//	chs[i], dsc[i] = text.TitleAndContent(option)
	//}

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
		Prompt: msg,
	})

	if err != nil {
		if errors.Is(err, tea.ErrInterrupted) {
			return "", err
		}
		return "", err
	}

	return model.Selected, nil
}

func Sendf(chs []search.Choice, format string, vars ...any) (string, error) {
	return Send(chs, fmt.Sprintf(format, vars...))
}
