package search

import (
	"github.com/razshare/frizzante/v2/tui/viewport"
)

func Reset(search *Search, viewport *viewport.Viewport) {
	search.Value = ""
	search.Active = false
	search.Filtered = search.Choices
	viewport.Cursor = 0
	viewport.Offset = 0
}
