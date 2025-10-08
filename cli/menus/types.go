package menus

import (
	"github.com/razshare/frizzante/tui/search"
)

type Menu struct {
	Items []Item
}

type Item struct {
	Choice  search.Choice
	Active  func() bool // Active returns true if the user has typed the choice directly in the terminal.
	Handler func() error
	Hidden  bool
}
