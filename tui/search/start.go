package search

import (
	"github.com/razshare/frizzante/tui/viewport"
	"strings"
)

func Filter(search *Search, vport *viewport.Viewport) {
	input := strings.ToLower(search.Input.Value())
	if input == "" {
		search.Filtered = search.Choices
	} else {
		filtered := make([]string, 0, len(search.Choices)/2)
		for _, choice := range search.Choices {
			if strings.Contains(strings.ToLower(choice), input) {
				filtered = append(filtered, choice)
			}
		}
		search.Filtered = filtered
	}
	vport.Cursor = 0
	vport.Start = 0
}
