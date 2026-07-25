package select_npm_packages

import (
	"time"

	"github.com/razshare/frizzante/v2/tui/program"
	"github.com/razshare/frizzante/v2/tui/search"
	"github.com/razshare/frizzante/v2/tui/viewport"
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
