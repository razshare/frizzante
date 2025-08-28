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

func Send(choices []search.Choice, message string) (selected string, err error) {
	input := textinput.New()
	input.Width = 80
	var model *Model
	model, err = program.Run(&Model{
		Prompt:   message,
		Viewport: &viewport.Viewport{Visible: 6},
		Search: &search.Search{
			Choices:  choices,
			Filtered: choices,
			Input:    input,
		},
	})

	if err != nil {
		if errors.Is(err, tea.ErrInterrupted) {
			return
		}
		return
	}

	selected = model.Selected

	return
}

func Sendf(choices []search.Choice, format string, vars ...any) (selected string, err error) {
	return Send(choices, fmt.Sprintf(format, vars...))
}
