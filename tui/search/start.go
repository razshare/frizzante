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
		filtered := make([]Choice, 0)
		for _, choice := range search.Choices {
			if strings.Contains(strings.ToLower(choice.Id), input) {
				filtered = append(filtered, choice)
			}
		}
		search.Filtered = filtered
	}
	vport.Cursor = 0
	vport.Start = 0
}
