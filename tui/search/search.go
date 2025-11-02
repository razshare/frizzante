package search

import "github.com/charmbracelet/bubbles/textinput"

type Search struct {
	Input    textinput.Model // Search input field
	Choices  []Choice        // All available choices
	Filtered []Choice        // Choices after filtering
	Active   bool            // Whether we're in search mode
}
