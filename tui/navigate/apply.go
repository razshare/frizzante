package navigate

import (
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/viewport"
)

func Apply(search *search.Search, vport *viewport.Viewport, direction int) {
	count := len(search.Filtered)
	if count == 0 {
		return
	}

	vport.Cursor = (vport.Cursor + direction + count) % count

	if vport.Cursor < vport.Start {
		vport.Start = vport.Cursor
	} else if vport.Cursor >= vport.Start+vport.Visible {
		vport.Start = vport.Cursor - vport.Visible + 1
	}
}
