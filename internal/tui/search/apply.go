package search

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/internal/tui/viewport"
)

func Apply(search *Search, viewport *viewport.Viewport, message tea.KeyMsg) tea.Cmd {
	viewport.Cursor = 0
	previous := search.Input.Value()
	var cmd tea.Cmd
	search.Input, cmd = search.Input.Update(message)

	if current := search.Input.Value(); current != previous {
		if current == "" {
			Reset(search, viewport)
		} else {
			Filter(search, viewport)
		}
	}

	return cmd
}
