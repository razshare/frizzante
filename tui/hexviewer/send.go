package hexviewer

import (
	"strings"

	"github.com/razshare/frizzante/tui/program"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/viewport"
)

func Send(title string, text string) (err error) {
	lines := strings.Split(text, "\n")
	choices := make([]search.Choice, len(lines))

	for index, line := range lines {
		choices[index] = search.Choice{
			Id:          line,
			Description: line,
		}
	}

	_, err = program.Run(&Model{
		Prompt:   title,
		Viewport: &viewport.Viewport{Visible: 12},
		Search: &search.Search{
			Choices:  choices,
			Filtered: choices,
		},
	})

	return
}
