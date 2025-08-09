package singleselect

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/razshare/frizzante/tui/program"
)

func Send(prompt string, options []string) (string, error) {
	// Initialize the search input
	searchInput := textinput.New()
	searchInput.Width = 80
	result, err := program.Run(&Model{
		Choices:         options,
		FilteredChoices: options,
		Prompt:          prompt,
		SearchInput:     searchInput,
		MaxVisible:      6,
		ViewportStart:   0,
		Cursor:          0,
		Searching:       false,
	})
	if err != nil {
		return "", err
	}
	if result.Selected == "" {
		return "", fmt.Errorf("no selection made")
	}
	return result.Selected, nil
}
