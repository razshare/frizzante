package menu

import "github.com/razshare/frizzante/tui/search"

type Menu struct {
	Items []Item
}

type Item struct {
	Choice  search.Choice
	Inlined func() bool // Inlined returns true if the user has typed the choice directly in the terminal.
	Handler func() error
}
