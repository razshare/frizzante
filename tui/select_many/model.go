package select_many

import (
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/viewport"
)

type Model struct {
	Prompt   string
	Selected []string
	Search   *search.Search
	Viewport *viewport.Viewport
}
