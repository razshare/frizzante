package multiselect

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/razshare/frizzante/stack"
	"github.com/razshare/frizzante/tui/program"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/text"
	"github.com/razshare/frizzante/tui/viewport"
)

func Send(opts []string, msg string) ([]string, error) {
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
		Selected: make([]string, 0),
		Prompt:   msg,
	})

	if err != nil {
		return nil, errors.New(err.Error() + "\n" + stack.Trace())
	}

	return model.Selected, nil
}

func Sendf(opts []string, format string, vars ...any) []string {
	return Sendf(opts, fmt.Sprintf(format, vars...))
}
