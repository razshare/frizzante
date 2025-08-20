package npmselect

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/razshare/frizzante/tui/program"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/viewport"
	"strings"
	"time"
)

func Send() ([]string, error) {
	in := textinput.New()
	in.Width = 80
	m, err := program.Run(&Model{
		Prompt:    "search npm packages",
		Viewport:  &viewport.Viewport{Visible: 6},
		Selected:  make([]string, 0),
		Debounce:  time.Second,
		Debouncer: time.NewTimer(time.Second),
		Search: &search.Search{
			Choices:  []search.Choice{},
			Filtered: []search.Choice{},
			Input:    in,
		},
	})

	if err != nil {
		return nil, err
	}

	ns := make([]string, 0, len(m.Selected))
	for _, id := range m.Selected {
		// Remove version suffix if present (e.g., "package@1.0.0" -> "package")
		if idx := strings.IndexByte(id, '@'); idx != -1 {
			ns = append(ns, id[:idx])
		} else {
			ns = append(ns, id)
		}
	}

	return ns, nil
}
