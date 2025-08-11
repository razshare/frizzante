package singleselect

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/program"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/viewport"
)

func Send(options []string, prompt string) string {
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
		messages.Fatal(err)
	}

	return result.Selected
}

func Sendf(options []string, format string, vars ...any) string {
	return Send(options, fmt.Sprintf(format, vars...))
}
