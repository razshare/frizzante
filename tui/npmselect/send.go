package npmselect

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/razshare/frizzante/tui/program"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/viewport"
	"strings"
	"time"
)

func Send() (selected []string, err error) {
	input := textinput.New()
	input.Width = 80
	var model *Model
	model, err = program.Run(&Model{
		Prompt:    "search npm packages",
		Viewport:  &viewport.Viewport{Visible: 6},
		Selected:  make([]string, 0),
		Debounce:  time.Second,
		Debouncer: time.NewTimer(time.Second),
		Search: &search.Search{
			Choices:  []search.Choice{},
			Filtered: []search.Choice{},
			Input:    input,
		},
	})

	if err != nil {
		return
	}

	selected = make([]string, 0, len(model.Selected))
	for _, id := range model.Selected {
		// Remove version suffix if present (e.g., "package@1.0.0" -> "package")
		if index := strings.IndexByte(id, '@'); index != -1 {
			selected = append(selected, id[:index])
		} else {
			selected = append(selected, id)
		}
	}

	return
}
