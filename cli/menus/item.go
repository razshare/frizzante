package menus

import "github.com/razshare/frizzante/tui/search"

type Item struct {
	Choice  search.Choice
	Ids     []string
	Handler func(value string) error
	Hidden  bool
}
