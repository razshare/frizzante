package hexviewer

import (
	"github.com/razshare/frizzante/v2/tui/search"
	"github.com/razshare/frizzante/v2/tui/viewport"
)

// Model defines text viewer options
type Model struct {
	Prompt   string
	Search   *search.Search
	Viewport *viewport.Viewport
}
