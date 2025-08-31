package multiselect

import (
	"github.com/razshare/frizzante/internal/tui/search"
	"github.com/razshare/frizzante/internal/tui/viewport"
)

type Model struct {
	Selected []string
	Prompt   string
	Search   *search.Search
	Viewport *viewport.Viewport
}
