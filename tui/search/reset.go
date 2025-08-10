package search

import "github.com/razshare/frizzante/tui/viewport"

func Reset(search *Search, vport *viewport.Viewport) {
	search.Input.SetValue("")
	search.Active = false
	search.Filtered = search.Choices
	vport.Cursor = 0
	vport.Start = 0
}
