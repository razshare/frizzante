package hexviewer

import (
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/viewport"
)

// Model defines text viewer options
type Model struct {
	Prompt   string
	Search   *search.Search
	Viewport *viewport.Viewport
}
