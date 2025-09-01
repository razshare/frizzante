package singleselect

import (
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/viewport"
)

// Model defines single selection options
type Model struct {
	Prompt   string
	Selected string
	Search   *search.Search
	Viewport *viewport.Viewport
}
