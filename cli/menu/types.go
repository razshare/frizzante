package menu

import "github.com/razshare/frizzante/tui/search"

type Menu struct {
	Items []Item
}

type Item struct {
	Choice  search.Choice
	Inlined func() bool
	Handler func() error
}
