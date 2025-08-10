package search

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/tui/viewport"
)

func Apply(search *Search, vport *viewport.Viewport, msg tea.KeyMsg) tea.Cmd {
	preval := search.Input.Value()
	var cmd tea.Cmd
	search.Input, cmd = search.Input.Update(msg)

	if newval := search.Input.Value(); newval != preval {
		if newval == "" {
			Reset(search, vport)
		} else {
			Filter(search, vport)
		}
	}
	return cmd
}
