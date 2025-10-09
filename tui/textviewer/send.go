package textviewer

import (
	"errors"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/tui/program"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/viewport"
)

func Send(title string, text string) (err error) {
	lines := strings.Split(text, "\n")
	choices := make([]search.Choice, len(lines))

	for index, line := range lines {
		choices[index] = search.Choice{
			Id:          line,
			Description: line,
		}
	}

	input := textinput.New()
	input.Width = 80
	_, err = program.Run(&Model{
		Prompt:   title,
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

	return
}
