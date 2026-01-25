package select_npm_packages

import (
	"time"

	"github.com/razshare/frizzante/tui/program"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/viewport"
)

func Send() (selected []string, err error) {
	var model *Model
	if model, err = program.Run(&Model{
		Prompt:    "search npm packages",
		Viewport:  &viewport.Viewport{Visible: 6},
		Selected:  make([]string, 0),
		Debounce:  time.Second,
		Debouncer: time.NewTimer(time.Second),
		Search: &search.Search{
			Choices:  []search.Choice{},
			Filtered: []search.Choice{},
		},
	}); err != nil {
		return
	}
	selected = model.Selected
	return
}
