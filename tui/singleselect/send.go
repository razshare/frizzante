package singleselect

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/program"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/text"
	"github.com/razshare/frizzante/tui/viewport"
)

func Send(opts []string, msg string) string {
	c := len(opts)
	chs := make([]string, c)
	dsc := make([]string, c)

	for i, option := range opts {
		chs[i], dsc[i] = text.TitleAndContent(option)
	}

	// Initialize the search input
	searchInput := textinput.New()
	searchInput.Width = 80
	model, err := program.Run(&Model{
		Search: &search.Search{
			Active:       false,
			Choices:      chs,
			Filtered:     chs,
			Descriptions: dsc,
			Input:        searchInput,
		},
		Viewport: &viewport.Viewport{
			Visible: 6,
			Start:   0,
			Cursor:  0,
		},
		Prompt: msg,
	})

	if err != nil {
		messages.Fatal(err)
	}

	return model.Selected
}

func Sendf(options []string, format string, vars ...any) string {
	return Send(options, fmt.Sprintf(format, vars...))
}
