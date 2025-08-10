package singleselect

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/razshare/frizzante/tui/program"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/viewport"
)

func Send(prompt string, options []string) (string, error) {
	// Initialize the search input
	searchInput := textinput.New()
	searchInput.Width = 80
	result, err := program.Run(&Model{
		Search: &search.Search{
			Active:   false,
			Choices:  options,
			Filtered: options,
			Input:    searchInput,
		},
		Viewport: &viewport.Viewport{
			Visible: 6,
			Start:   0,
			Cursor:  0,
		},
		Prompt: prompt,
	})

	if err != nil {
		return "", err
	}

	return result.Selected, nil
}
