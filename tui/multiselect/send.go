package multiselect

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/program"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/viewport"
)

func Send(options []string, prompt string) []string {
	// Initialize the search input
	searchInput := textinput.New()
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
		Selected: make(map[int]bool),
		Prompt:   prompt,
	})

	if err != nil {
		messages.Fatal(err)
	}

	var selections []string
	for i, choice := range result.Search.Choices {
		if result.Selected[i] {
			selections = append(selections, choice)
		}
	}

	return selections
}

func Sendf(options []string, format string, vars ...any) []string {
	return Sendf(options, fmt.Sprintf(format, vars...))
}
