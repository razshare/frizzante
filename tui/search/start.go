package search

import (
	"github.com/razshare/frizzante/tui/viewport"
	"strings"
)

func Filter(search *Search, viewport *viewport.Viewport) {
	input := strings.ToLower(search.Input.Value())
	if input == "" {
		search.Filtered = search.Choices
	} else {
		filtered := make([]Choice, 0)
		for _, choice := range search.Choices {
			if strings.Contains(strings.ToLower(choice.Id), input) {
				filtered = append(filtered, choice)
			}
		}
		search.Filtered = filtered
	}
	viewport.Cursor = 0
	viewport.Start = 0
}
