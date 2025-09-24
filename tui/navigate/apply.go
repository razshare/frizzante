package navigate

import (
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/viewport"
)

func Apply(search *search.Search, viewport *viewport.Viewport, direction int) {
	count := len(search.Filtered)
	if count == 0 {
		return
	}

	viewport.Cursor = ((viewport.Cursor+direction)%count + count) % count

	if viewport.Cursor < viewport.Offset {
		viewport.Offset = viewport.Cursor
	} else if viewport.Cursor >= viewport.Offset+viewport.Visible {
		viewport.Offset = viewport.Cursor - viewport.Visible + 1
	}
}
