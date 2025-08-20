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
	input := textinput.New()
	input.Width = 80
	model, err := program.Run(&Model{
		Prompt:   msg,
		Viewport: &viewport.Viewport{Visible: 6},
		Search: &search.Search{
			Choices:  chs,
			Filtered: chs,
			Input:    input,
		},
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
