package menus

import "github.com/razshare/frizzante/tui/search"

type Item struct {
	Choice  search.Choice
	Active  ActivationFunction
	Handler ActionFunction
	Hidden  bool
}
