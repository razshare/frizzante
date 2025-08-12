package multiselect

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/program"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/viewport"
	"strings"
)

func Send(opts []string, msg string) []string {
	var sb strings.Builder
	c := len(opts)
	chs := make([]string, c)
	dsc := make([]string, c)

	for i, option := range opts {
		p := strings.SplitN(strings.TrimSpace(option), "\n", 2)

		if len(p) > 1 {
			for _, l := range strings.Split(p[1], "\n") {
				trm := strings.TrimSpace(l)
				if trm == "" {
					continue
				}
				sb.WriteString(trm + " ")
			}
			dsc[i] = sb.String()
			sb.Reset()
		} else {
			dsc[i] = ""
		}

		chs[i] = p[0]
	}

	// Initialize the search input
	searchInput := textinput.New()
	result, err := program.Run(&Model{
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
		Selected: make(map[int]bool),
		Prompt:   msg,
	})

	if err != nil {
		messages.Fatal(err)
	}

	var selections []string
	for i, choice := range result.Search.Choices {
		if result.Selected[i] {
			selections = append(selections, choice)
		}
	}

	return selections
}

func Sendf(opts []string, format string, vars ...any) []string {
	return Sendf(opts, fmt.Sprintf(format, vars...))
}
