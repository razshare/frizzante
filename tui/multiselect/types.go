package multiselect

import (
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/viewport"
)

type Model struct {
	Selected []string
	Prompt   string
	Search   *search.Search
	Viewport *viewport.Viewport
}
