package search

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/tui/viewport"
)

func Apply(search *Search, viewport *viewport.Viewport, message tea.KeyMsg) tea.Cmd {
	search.Active = true
	viewport.Cursor = 0
	previous := search.Value
	var cmd tea.Cmd

	// vscode has a weird bug where it will send a "ctrl+w" whenever the user presses "backspace" in the integrated terminal,
	// so we're including tea.KeyCtrlW to try to fix that for the user.
	// https://stackoverflow.com/questions/52806758/visual-studio-code-ctrlbackspace-not-working-in-integrated-terminal
	if message.Type == tea.KeyBackspace || message.Type == tea.KeyCtrlH || message.Type == tea.KeyCtrlW {
		length := len(search.Value)
		if length > 0 {
			search.Value = search.Value[:length-1]
		}
	} else {
		var content string
		if content = message.String(); len(content) == 1 {
			// we only accept single characters
			search.Value += content
		}
	}

	if current := search.Value; current != previous {
		if current == "" {
			Reset(search, viewport)
		} else {
			Filter(search, viewport)
		}
	}

	return cmd
}
