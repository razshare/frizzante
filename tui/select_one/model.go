package select_one

import (
	"github.com/razshare/frizzante/v2/tui/search"
	"github.com/razshare/frizzante/v2/tui/viewport"
)

// Model defines single selection options
type Model struct {
	Prompt   string
	Selected string
	Search   *search.Search
	Viewport *viewport.Viewport
}
