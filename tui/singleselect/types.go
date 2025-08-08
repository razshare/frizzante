package singleselect

import "github.com/charmbracelet/bubbles/textinput"

// Model defines single selection options
type Model struct {
	Choices         []string        // All available Choices
	FilteredChoices []string        // Choices after filtering
	Cursor          int             // Current Cursor position in filtered list
	Prompt          string          // The Prompt to display
	Selected        string          // The Selected choice
	SearchInput     textinput.Model // Search input field
	Searching       bool            // Whether we're in search mode
	ViewportStart   int             // Start index for viewport (for scrolling)
	MaxVisible      int             // Maximum visible items (5 by default)
}
