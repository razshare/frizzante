package select_many

import (
	"github.com/razshare/frizzante/v2/tui/search"
	"github.com/razshare/frizzante/v2/tui/viewport"
)

type Model struct {
	Prompt   string
	Selected []string
	Search   *search.Search
	Viewport *viewport.Viewport
}
