package multiselect

import (
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/viewport"
)

type Model struct {
	Selected map[int]bool
	Prompt   string
	Search   *search.Search
	Viewport *viewport.Viewport
}
