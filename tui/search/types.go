package search

import "github.com/charmbracelet/bubbles/textinput"

type Search struct {
	Active   bool            // Whether we're in search mode
	Choices  []Choice        // All available choices
	Filtered []Choice        // Choices after filtering
	Input    textinput.Model // Search input field
}

type Choice struct {
	Id          string
	Description string
}
