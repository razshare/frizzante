package singleselect

import (
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/viewport"
)

// Model defines single selection options
type Model struct {
	Prompt        string // The Prompt to display
	Selected      string // The Selected choice
	Search        *search.Search
	Viewport      *viewport.Viewport
	SoftInterrupt bool
}
