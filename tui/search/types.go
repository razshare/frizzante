package search

import "github.com/charmbracelet/bubbles/textinput"

type Search struct {
	Active       bool            // Whether we're in search mode
	Choices      []string        // All available choices
	Descriptions []string        // All available descriptions
	Filtered     []string        // Choices after filtering
	Input        textinput.Model // Search input field
}
