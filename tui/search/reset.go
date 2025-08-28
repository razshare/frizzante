package search

import "github.com/razshare/frizzante/tui/viewport"

func Reset(search *Search, viewport *viewport.Viewport) {
	search.Input.SetValue("")
	search.Active = false
	search.Filtered = search.Choices
	viewport.Cursor = 0
	viewport.Start = 0
}
